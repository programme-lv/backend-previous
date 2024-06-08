package submissions

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sort"

	set "github.com/deckarep/golang-set/v2"
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/postgres/proglv/public/table"
	"github.com/programme-lv/director/msg"
	"google.golang.org/grpc/metadata"
)

type TestingDirectorConn struct {
	GRPCClient msg.DirectorClient
	Password   string
}

type TestDownloadURLProvider interface {
	GetTestDownloadURL(testSHA256 string) (string, error)
}

func EvaluateSubmission(db qrm.DB,
	submissionID int64, taskVersionID int64,
	urlProvider TestDownloadURLProvider,
	directorConn TestingDirectorConn) error {

	evalID, err := insertNewEvaluation(db, taskVersionID)
	if err != nil {
		return err
	}

	// set this evaluation as the visible one
	stmtUpdateEval := table.TaskSubmissions.UPDATE(
		table.TaskSubmissions.VisibleEvalID).
		SET(postgres.Int64(evalID)).
		WHERE(table.TaskSubmissions.ID.EQ(
			postgres.Int64(submissionID)))
	_, err = stmtUpdateEval.Exec(db)
	if err != nil {
		return err
	}

	req := msg.EvaluationRequest{
		Submission: "",
		Language: &msg.ProgrammingLanguage{
			Id:               "",
			Name:             "",
			CodeFilename:     "",
			CompileCmd:       new(string),
			CompiledFilename: new(string),
			ExecuteCmd:       "",
		},
		Limits: &msg.RuntimeLimits{
			CPUTimeMillis: 0,
			MemKibiBytes:  0,
		},
		EvalType:       "",
		Tests:          []*msg.Test{},
		TestlibChecker: "",
	}

	err = populateSubmissionAndLangID(db, submissionID, &req)
	if err != nil {
		return err
	}

	err = populateProgrammingLanguage(db, req.Language.Id, &req)
	if err != nil {
		return err
	}

	taskVersion, err := getTaskVersion(db, taskVersionID)
	if err != nil {
		return err
	}

	req.EvalType = taskVersion.TestingTypeID
	req.Limits.CPUTimeMillis = taskVersion.TimeLimMs
	req.Limits.MemKibiBytes = taskVersion.MemLimKibibytes

	if taskVersion.CheckerID != nil {
		checkerID := *taskVersion.CheckerID
		err = populateTestlibChecker(db, checkerID, &req)
		if err != nil {
			return nil
		}
	}

	err = populateTests(db, taskVersionID, urlProvider, &req)
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"authorization": directorConn.Password})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	evaluationClient, err := directorConn.GRPCClient.EvaluateSubmission(ctx, &req)
	if err != nil {
		return err
	}

	go func() {
		fbProc := NewEvalFeedbackProcessor(db, evalID)
		for {
			res, err := evaluationClient.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				slog.Error("error receiving evaluation feedback", err)
				break
			}

			err = fbProc.Process(res)
			if err != nil {
				slog.Error("error processing evaluation feedback", err)
				break
			}
		}
	}()

	return nil
}

func insertNewEvaluation(db qrm.DB, taskVersionID int64) (int64, error) {
	evaluation := model.Evaluations{
		EvalStatusID:  "IQ",
		TaskVersionID: taskVersionID,
	}

	insertStmt := table.Evaluations.INSERT(
		table.Evaluations.EvalStatusID,
		table.Evaluations.TaskVersionID,
	).MODEL(evaluation).RETURNING(table.Evaluations.ID)
	err := insertStmt.Query(db, &evaluation)
	if err != nil {
		return 0, err
	}
	return evaluation.ID, nil
}

func populateSubmissionAndLangID(db qrm.DB, submID int64, req *msg.EvaluationRequest) error {
	stmtSelSubm := postgres.SELECT(
		table.TaskSubmissions.ProgrammingLangID,
		table.TaskSubmissions.Submission,
	).FROM(table.TaskSubmissions).
		WHERE(table.TaskSubmissions.ID.EQ(
			postgres.Int64(submID)))

	var subm model.TaskSubmissions
	err := stmtSelSubm.Query(db, &subm)
	if err != nil {
		return err
	}

	req.Submission = subm.Submission
	req.Language.Id = subm.ProgrammingLangID
	return nil
}

func populateProgrammingLanguage(db qrm.DB, programmingLangID string, req *msg.EvaluationRequest) error {
	stmtSelLang := postgres.SELECT(
		table.ProgrammingLanguages.AllColumns,
	).FROM(table.ProgrammingLanguages).
		WHERE(table.ProgrammingLanguages.ID.EQ(
			postgres.String(programmingLangID)))

	var lang model.ProgrammingLanguages
	err := stmtSelLang.Query(db, &lang)
	if err != nil {
		return err
	}

	req.Language.Name = lang.FullName
	req.Language.CodeFilename = lang.CodeFilename
	req.Language.CompileCmd = lang.CompileCmd
	req.Language.CompiledFilename = lang.CompiledFilename
	req.Language.ExecuteCmd = lang.ExecuteCmd
	return nil
}

func getTaskVersion(db qrm.DB, taskVersionID int64) (*model.TaskVersions, error) {
	stmtSelTaskVersion := postgres.SELECT(
		table.TaskVersions.AllColumns,
	).FROM(table.TaskVersions).
		WHERE(table.TaskVersions.ID.EQ(
			postgres.Int64(taskVersionID)))

	var taskVersion model.TaskVersions
	err := stmtSelTaskVersion.Query(db, &taskVersion)
	if err != nil {
		return nil, err
	}
	return &taskVersion, nil
}

func populateTestlibChecker(db qrm.DB, checkerID int64, req *msg.EvaluationRequest) error {
	stmtSelChecker := postgres.SELECT(
		table.TestlibCheckers.AllColumns,
	).FROM(table.TestlibCheckers).
		WHERE(table.TestlibCheckers.ID.EQ(
			postgres.Int64(checkerID)))

	var checker model.TestlibCheckers
	err := stmtSelChecker.Query(db, &checker)
	if err != nil {
		return err
	}
	req.TestlibChecker = checker.Code
	return nil
}

func populateTests(db qrm.DB, taskVersionID int64, urls TestDownloadURLProvider, req *msg.EvaluationRequest) error {
	selectTestSetIDStmt := postgres.SELECT(
		table.TaskVersions.TestSetID,
	).FROM(table.TaskVersions).
		WHERE(table.TaskVersions.ID.EQ(
			postgres.Int64(taskVersionID)))

	var taskVersionRecord model.TaskVersions
	err := selectTestSetIDStmt.Query(db, &taskVersionRecord)
	if err != nil {
		return err
	}

	if taskVersionRecord.TestSetID == nil {
		return nil
	}

	testSetID := *taskVersionRecord.TestSetID

	stmtSelectTestInputs := postgres.SELECT(
		table.TestSetTests.ID,
		table.TextFiles.Sha256,
		table.TextFiles.Content,
		table.TextFiles.Compression,
	).FROM(table.TestSetTests.
		LEFT_JOIN(table.TextFiles, table.TestSetTests.InputTextFileID.EQ(table.TextFiles.ID))).
		WHERE(table.TestSetTests.TestSetID.EQ(
			postgres.Int64(testSetID)))

	var inputs []struct {
		model.TestSetTests
		model.TextFiles
	}
	err = stmtSelectTestInputs.Query(db, &inputs)
	if err != nil {
		return err
	}

	stmtSelectTestAnswers := postgres.SELECT(
		table.TestSetTests.ID,
		table.TextFiles.Sha256,
		table.TextFiles.Content,
		table.TextFiles.Compression,
	).FROM(table.TestSetTests.
		LEFT_JOIN(table.TextFiles, table.TestSetTests.AnswerTextFileID.EQ(table.TextFiles.ID))).
		WHERE(table.TestSetTests.TestSetID.EQ(
			postgres.Int64(testSetID)))
	var answers []struct {
		model.TestSetTests
		model.TextFiles
	}
	err = stmtSelectTestAnswers.Query(db, &answers)
	if err != nil {
		return err
	}

	type testPart struct {
		id          int64
		sha256      string
		content     *string
		compression string
	}

	testIdSet := set.NewSet[int64]()
	idInputMap := make(map[int64]testPart)
	idAnswerMap := make(map[int64]testPart)

	for _, input := range inputs {
		idInputMap[input.TestSetTests.ID] = testPart{
			id:          input.TestSetTests.ID,
			sha256:      input.Sha256,
			content:     input.TextFiles.Content,
			compression: input.Compression,
		}
		testIdSet.Add(input.TestSetTests.ID)
	}

	for _, answer := range answers {
		idAnswerMap[answer.TestSetTests.ID] = testPart{
			id:          answer.TestSetTests.ID,
			sha256:      answer.Sha256,
			content:     answer.TextFiles.Content,
			compression: answer.Compression,
		}
		testIdSet.Add(answer.TestSetTests.ID)
	}

	for _, testID := range testIdSet.ToSlice() {
		input, okInp := idInputMap[testID]
		if !okInp {
			return fmt.Errorf("input test not found")
		}
		answer, okAns := idAnswerMap[testID]
		if !okAns {
			return fmt.Errorf("answer test not found")
		}

		test := &msg.Test{
			Id:             testID,
			InSha256:       input.sha256,
			InDownloadUrl:  nil,
			InContent:      input.content,
			AnsSha256:      answer.sha256,
			AnsDownloadUrl: nil,
			AnsContent:     answer.content,
		}
		if test.InContent == nil {
			url, err := urls.GetTestDownloadURL(input.sha256)
			if err != nil {
				return err
			}
			test.InDownloadUrl = &url
		}
		if test.AnsContent == nil {
			url, err := urls.GetTestDownloadURL(answer.sha256)
			if err != nil {
				return err
			}
			test.AnsDownloadUrl = &url
		}
		req.Tests = append(req.Tests, test)
	}

	sort.Slice(req.Tests, func(i, j int) bool {
		return req.Tests[i].Id < req.Tests[j].Id
	})

	return nil
}
