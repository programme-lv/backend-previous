package submissions

import (
	"log/slog"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
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
	case *pb.EvaluationFeedback_FinishWithCompilationError:
		slog.Debug("received \"FinishWithCompilationError\" feedback", "body", res.GetFinishWithCompilationError())
		stmt := table.Evaluations.UPDATE(table.Evaluations.EvalStatusID).
			SET(postgres.String("CE")).
			WHERE(table.Evaluations.ID.EQ(postgres.Int64(fb.evalID)))
		_, err := stmt.Exec(fb.db)
		return err
	case *pb.EvaluationFeedback_StartTesting:
		// s.infoLog.Printf("StartTesting: %+v", res.GetStartTesting())
	case *pb.EvaluationFeedback_IgnoreTest:
		// s.infoLog.Printf("IgnoreTest: %+v", res.GetIgnoreTest())
	case *pb.EvaluationFeedback_StartTest:
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
