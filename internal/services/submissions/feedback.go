package submissions

import (
	"log/slog"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	pb "github.com/programme-lv/director/msg"
)

type EvalFeedbackProcessor struct {
	db     qrm.DB
	evalID int64
}

func NewEvalFeedbackProcessor(db qrm.DB, evaluationID int64) *EvalFeedbackProcessor {
	return &EvalFeedbackProcessor{
		db:     db,
		evalID: evaluationID,
	}
}

func (fb *EvalFeedbackProcessor) Process(res *pb.EvaluationFeedback) error {
	switch res.FeedbackTypes.(type) {
	case *pb.EvaluationFeedback_StartEvaluation:
		slog.Debug("received \"StartEvaluation\" feedback", "body", res.GetStartEvaluation())
		stmt := table.Evaluations.UPDATE(table.Evaluations.EvalStatusID).
			SET(postgres.String("R")).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err := stmt.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_FinishEvaluation:
		slog.Debug("received \"FinishEvaluation\" feedback", "body", res.GetFinishEvaluation())
		stmt := table.Evaluations.UPDATE(table.Evaluations.EvalStatusID).
			SET(postgres.String("F")).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err := stmt.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_FinishWithInernalServerError:
		slog.Debug("received \"FinishWithInernalServerError\" feedback", "body", res.GetFinishWithInernalServerError())
		slog.Error("evaluation finished with internal server error", "evaluation_id", fb.evalID, "error", res.GetFinishWithInernalServerError().ErrorMsg)
		stmt := table.Evaluations.UPDATE(table.Evaluations.EvalStatusID).
			SET(postgres.String("ISE")).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err := stmt.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_StartCompilation:
		slog.Debug("received \"StartCompilation\" feedback", "body", res.GetStartCompilation())
		stmt := table.Evaluations.UPDATE(table.Evaluations.EvalStatusID).
			SET(postgres.String("C")).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err := stmt.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_FinishCompilation:
		slog.Debug("received \"FinishCompilation\" feedback", "body", res.GetFinishCompilation())
		data := res.GetFinishCompilation().GetCompilationRData()
		// insert into runtime_data table (stdout,stderr,time_millis,memory_kibibytes,time_wall_millis,exit_code)
		stmt := table.RuntimeData.INSERT(
			table.RuntimeData.Stdout,
			table.RuntimeData.Stderr,
			table.RuntimeData.TimeMillis,
			table.RuntimeData.MemoryKibibytes,
			table.RuntimeData.TimeWallMillis,
			table.RuntimeData.ExitCode,
		).VALUES(
			data.Stdout,
			data.Stderr,
			data.CpuTimeMillis,
			data.MemKibiBytes,
			data.WallTimeMillis,
			data.ExitCode,
		).RETURNING(table.RuntimeData.ID)
		var cRunData model.RuntimeData
		err := stmt.Query(fb.db, &cRunData)
		if err != nil {
			return err
		}
		// link runtime_data to evaluation
		stmt2 := table.Evaluations.UPDATE(table.Evaluations.CompilationDataID).
			SET(postgres.Int64(cRunData.ID)).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err = stmt2.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_FinishWithCompilationError:
		slog.Debug("received \"FinishWithCompilationError\" feedback", "body", res.GetFinishWithCompilationError())
		stmt := table.Evaluations.UPDATE(table.Evaluations.EvalStatusID).
			SET(postgres.String("CE")).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err := stmt.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_StartTesting:
		slog.Debug("received \"StartTesting\" feedback", "body", res.GetStartTesting())
		stmt := table.Evaluations.UPDATE(table.Evaluations.EvalStatusID).
			SET(postgres.String("T")).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err := stmt.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_IgnoreTest:
		slog.Debug("received \"IgnoreTest\" feedback", "body", res.GetIgnoreTest())
		testID := res.GetIgnoreTest().TestId
		stmt := table.EvaluationTestResults.UPDATE(
			table.EvaluationTestResults.EvalStatusID,
		).SET(postgres.String("IG")).
			WHERE(table.EvaluationTestResults.EvaluationID.EQ(postgres.Int64(fb.evalID)).
				AND(table.EvaluationTestResults.TaskVTestID.EQ(postgres.Int64(testID))))
		_, err := stmt.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_StartTest:
		slog.Debug("received \"StartTest\" feedback", "body", res.GetStartTest())
		testID := res.GetStartTest().TestId
		stmt := table.EvaluationTestResults.INSERT(
			table.EvaluationTestResults.EvaluationID,
			table.EvaluationTestResults.EvalStatusID,
			table.EvaluationTestResults.TaskVTestID,
		).VALUES(postgres.Int64(fb.evalID), postgres.String("T"), postgres.Int64(testID))
		_, err := stmt.Exec(fb.db)
		return err
		// s.infoLog.Printf("StartTest: %+v", re,s.GetStartTest())
	case *pb.EvaluationFeedback_ReportTestSubmissionRuntimeData:
		// s.infoLog.Printf("ReportTestSubmissionRuntimeData: %+v", res.GetReportTestSubmissionRuntimeData())
	case *pb.EvaluationFeedback_FinishTestWithLimitExceeded:
		// s.infoLog.Printf("FinishTestWithLimitExceeded: %+v", res.GetFinishTestWithLimitExceeded())
	case *pb.EvaluationFeedback_FinishTestWithRuntimeError:
		// s.infoLog.Printf("FinishTestWithRuntimeError: %+v", res.GetFinishTestWithRuntimeError())
	case *pb.EvaluationFeedback_ReportTestCheckerRuntimeData:
		// s.infoLog.Printf("ReportTestCheckerRuntimeData: %+v", res.GetReportTestCheckerRuntimeData())
	case *pb.EvaluationFeedback_FinishTestWithVerdictAccepted:
		// s.infoLog.Printf("FinishTestWithVerdictAccepted: %+v", res.GetFinishTestWithVerdictAccepted())
	case *pb.EvaluationFeedback_FinishTestWithVerdictWrongAnswer:
		// s.infoLog.Printf("FinishTestWithVerdictWrongAnswer: %+v", res.GetFinishTestWithVerdictWrongAnswer())
	case *pb.EvaluationFeedback_IncrementScore:
		// s.infoLog.Printf("IncrementScore: %+v", res.GetIncrementScore())
	}
	return nil
}
