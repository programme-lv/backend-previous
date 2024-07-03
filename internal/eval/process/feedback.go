package evalsubm

import (
	"log/slog"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/table"
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
		// FINISH
	case *pb.EvaluationFeedback_FinishWithInernalServerError:
		slog.Debug("received \"FinishWithInernalServerError\" feedback", "body", res.GetFinishWithInernalServerError())
		slog.Error("evaluation finished with internal server error", "evaluation_id", fb.evalID, "error", res.GetFinishWithInernalServerError().ErrorMsg)
		stmt := table.Evaluations.UPDATE(table.Evaluations.EvalStatusID).
			SET(postgres.String("ISE")).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err := stmt.Exec(fb.db)
		return err
		//FINISH
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
		// FINISH
	case *pb.EvaluationFeedback_StartTesting:
		slog.Debug("received \"StartTesting\" feedback", "body", res.GetStartTesting())
		body := res.GetStartTesting()
		maxScore := body.GetMaxScore()
		stmt := table.Evaluations.UPDATE(
			table.Evaluations.EvalStatusID,
			table.Evaluations.EvalPossibleScore,
		).
			SET(postgres.String("T"), postgres.Int64(maxScore)).
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
	case *pb.EvaluationFeedback_ReportTestSubmissionRuntimeData:
		slog.Debug("received \"ReportTestSubmissionRuntimeData\" feedback", "body", res.GetReportTestSubmissionRuntimeData())
		data := res.GetReportTestSubmissionRuntimeData().GetRData()
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
		var rRunData model.RuntimeData
		err := stmt.Query(fb.db, &rRunData)
		if err != nil {
			return err
		}
		stmt2 := table.EvaluationTestResults.UPDATE(table.EvaluationTestResults.ExecRDataID).
			SET(postgres.Int64(rRunData.ID)).
			WHERE(table.EvaluationTestResults.EvaluationID.EQ(postgres.Int64(fb.evalID)).
				AND(table.EvaluationTestResults.TaskVTestID.EQ(postgres.Int64(res.GetReportTestSubmissionRuntimeData().TestId))))
		_, err = stmt2.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_FinishTestWithLimitExceeded:
		slog.Debug("received \"FinishTestWithLimitExceeded\" feedback", "body", res.GetFinishTestWithLimitExceeded())
		body := res.GetFinishTestWithLimitExceeded()
		idlenessLimExc := body.GetIdlenessLimitExceeded()
		memoryLimExc := body.GetMemoryLimitExceeded()
		timeLimExc := body.GetIsCPUTimeExceeded()
		testID := res.GetFinishTestWithLimitExceeded().TestId

		var status string
		if idlenessLimExc {
			status = "ILE"
		} else if memoryLimExc {
			status = "MLE"
		} else if timeLimExc {
			status = "TLE"
		}

		stmt := table.EvaluationTestResults.UPDATE(
			table.EvaluationTestResults.EvalStatusID,
		).SET(
			postgres.String(status),
		).WHERE(
			table.EvaluationTestResults.EvaluationID.EQ(postgres.Int64(fb.evalID)).
				AND(table.EvaluationTestResults.TaskVTestID.EQ(postgres.Int64(testID))),
		)
		_, err := stmt.Exec(fb.db)
		return err
		// FINISH
	case *pb.EvaluationFeedback_FinishTestWithRuntimeError:
		slog.Debug("received \"FinishTestWithRuntimeError\" feedback", "body", res.GetFinishTestWithRuntimeError())
		testID := res.GetFinishTestWithRuntimeError().TestId
		stmt := table.EvaluationTestResults.UPDATE(
			table.EvaluationTestResults.EvalStatusID,
		).SET(
			postgres.String("RE"),
		).WHERE(
			table.EvaluationTestResults.EvaluationID.EQ(postgres.Int64(fb.evalID)).
				AND(table.EvaluationTestResults.TaskVTestID.EQ(postgres.Int64(testID))),
		)
		_, err := stmt.Exec(fb.db)
		return err
		// FINISH
	case *pb.EvaluationFeedback_ReportTestCheckerRuntimeData:
		slog.Debug("received \"ReportTestCheckerRuntimeData\" feedback", "body", res.GetReportTestCheckerRuntimeData())
		body := res.GetReportTestCheckerRuntimeData()
		data := body.GetRData()
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
		stmt2 := table.EvaluationTestResults.UPDATE(table.EvaluationTestResults.CheckerRDataID).
			SET(postgres.Int64(cRunData.ID)).
			WHERE(table.EvaluationTestResults.EvaluationID.EQ(postgres.Int64(fb.evalID)).
				AND(table.EvaluationTestResults.TaskVTestID.EQ(postgres.Int64(body.TestId))))
		_, err = stmt2.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_FinishTestWithVerdictAccepted:
		slog.Debug("received \"FinishTestWithVerdictAccepted\" feedback", "body", res.GetFinishTestWithVerdictAccepted())
		testID := res.GetFinishTestWithVerdictAccepted().TestId
		stmt := table.EvaluationTestResults.UPDATE(
			table.EvaluationTestResults.EvalStatusID,
		).SET(
			postgres.String("AC"),
		).WHERE(
			table.EvaluationTestResults.EvaluationID.EQ(postgres.Int64(fb.evalID)).
				AND(table.EvaluationTestResults.TaskVTestID.EQ(postgres.Int64(testID))),
		)
		_, err := stmt.Exec(fb.db)
		return err
		// FINISH
	case *pb.EvaluationFeedback_FinishTestWithVerdictWrongAnswer:
		slog.Debug("received \"FinishTestWithVerdictWrongAnswer\" feedback", "body", res.GetFinishTestWithVerdictWrongAnswer())
		testID := res.GetFinishTestWithVerdictWrongAnswer().TestId
		stmt := table.EvaluationTestResults.UPDATE(
			table.EvaluationTestResults.EvalStatusID,
		).SET(
			postgres.String("WA"),
		).WHERE(
			table.EvaluationTestResults.EvaluationID.EQ(postgres.Int64(fb.evalID)).
				AND(table.EvaluationTestResults.TaskVTestID.EQ(postgres.Int64(testID))),
		)
		_, err := stmt.Exec(fb.db)
		return err
		// FINISH
	case *pb.EvaluationFeedback_IncrementScore:
		slog.Debug("received \"IncrementScore\" feedback", "body", res.GetIncrementScore())
		score := res.GetIncrementScore().GetDelta()
		stmt := table.Evaluations.UPDATE(
			table.Evaluations.EvalTotalScore,
		).SET(
			table.Evaluations.EvalTotalScore.ADD(postgres.Int64(score)),
		).WHERE(
			table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)),
		)
		_, err := stmt.Exec(fb.db)
		return err
	}
	return nil
}
