package submissions

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/domain"
)

func GetEvaluationObj(db qrm.DB, evalID int64, fillTests bool) (*domain.Evaluation, error) {
	res := domain.Evaluation{
		ID:             evalID,
		TaskVersionID:  0,
		StatusID:       "",
		ReceivedScore:  0,
		PossibleScore:  nil,
		CheckerRunData: nil,
		TestResults:    nil,
		CreatedAt:      time.Time{},
	}

	eval, err := selectEvaluationRecord(db, evalID)
	if err != nil {
		return nil, err
	}

	res.ID = eval.ID
	res.TaskVersionID = eval.TaskVersionID
	res.StatusID = eval.EvalStatusID
	res.ReceivedScore = eval.EvalTotalScore
	res.PossibleScore = eval.EvalPossibleScore
	res.CreatedAt = eval.CreatedAt

	if eval.CompilationDataID != nil {
		checkerRData, err := selectRuntimeDataRecord(db, *eval.CompilationDataID)
		if err != nil {
			return nil, err
		}

		res.CheckerRunData = &objects.RuntimeData{
			ID:              evalID,
			Stdout:          checkerRData.Stdout,
			Stderr:          checkerRData.Stderr,
			TimeMillis:      checkerRData.TimeWallMillis,
			MemoryKibibytes: checkerRData.MemoryKibibytes,
			TimeWallMillis:  checkerRData.TimeWallMillis,
			ExitCode:        checkerRData.ExitCode,
		}
	}

	if fillTests {
		testResults, err := selectEvalTestResWithRData(db, evalID)
		if err != nil {
			return nil, err
		}
		res.TestResults = make([]objects.EvalTestRes, len(testResults))
		for i, tr := range testResults {
			res.TestResults[i] = objects.EvalTestRes{
				ID:           tr.ETR.ID,
				EvaluationID: tr.ETR.EvaluationID,
				EvalStatusID: tr.ETR.EvalStatusID,
				TaskVTestID:  tr.ETR.TaskVTestID,
				ExecRData: &objects.RuntimeData{
					ID:              tr.RD1.ID,
					Stdout:          tr.RD1.Stdout,
					Stderr:          tr.RD1.Stderr,
					TimeMillis:      tr.RD1.TimeMillis,
					MemoryKibibytes: tr.RD1.MemoryKibibytes,
					TimeWallMillis:  tr.RD1.TimeWallMillis,
					ExitCode:        tr.RD1.ExitCode,
				},
				CheckerRData: &objects.RuntimeData{
					ID:              tr.RD2.ID,
					Stdout:          tr.RD2.Stdout,
					Stderr:          tr.RD2.Stderr,
					TimeMillis:      tr.RD2.TimeMillis,
					MemoryKibibytes: tr.RD2.MemoryKibibytes,
					TimeWallMillis:  tr.RD2.TimeWallMillis,
					ExitCode:        tr.RD2.ExitCode,
				},
			}
		}
	}

	return &res, nil
}

func selectEvaluationRecord(db qrm.DB, evalID int64) (*model.Evaluations, error) {
	stmt := postgres.SELECT(table.Evaluations.AllColumns).
		FROM(table.Evaluations).
		WHERE(table.Evaluations.ID.EQ(postgres.Int64(evalID)))

	var record model.Evaluations
	err := stmt.Query(db, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func selectRuntimeDataRecord(db qrm.DB, evalID int64) (*model.RuntimeData, error) {
	stmt := postgres.SELECT(table.RuntimeData.AllColumns).
		FROM(table.RuntimeData).
		WHERE(table.RuntimeData.ID.EQ(postgres.Int64(evalID)))

	var record model.RuntimeData
	err := stmt.Query(db, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

type EvaluationTestResultWithRuntimeData struct {
	ETR model.EvaluationTestResults `alias:"etr.*"`
	RD1 model.RuntimeData           `alias:"rd1.*"`
	RD2 model.RuntimeData           `alias:"rd2.*"`
}

func selectEvalTestResWithRData(db qrm.DB, evalID int64) ([]EvaluationTestResultWithRuntimeData, error) {
	etr := table.EvaluationTestResults.AS("etr")
	rd1 := table.RuntimeData.AS("rd1")
	rd2 := table.RuntimeData.AS("rd2")

	stmt := postgres.SELECT(etr.AllColumns, rd1.AllColumns, rd2.AllColumns).
		FROM(etr.
			INNER_JOIN(rd1, rd1.ID.EQ(etr.ExecRDataID)).
			INNER_JOIN(rd2, rd2.ID.EQ(etr.CheckerRDataID))).
		WHERE(etr.EvaluationID.EQ(postgres.Int64(evalID)))

	var records []EvaluationTestResultWithRuntimeData
	err := stmt.Query(db, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

/*
select
	etr.id as "etr_id",
	etr.evaluation_id as "etr_evaluation_id",
	etr.eval_status_id as "etr_eval_status_id",
	etr.task_v_test_id as "etr_task_v_test_id",
	r1.id as "r1_id", r1.stdout as "r1_stdout", r1.stderr as "r1_stderr",
	r1.time_millis as "r1_time_millis", r1.memory_kibibytes as "r1_memory_kibibytes",
	r1.time_wall_millis as "r1_time_wall_millis", r1.exit_code as "r1_exit_code",
	r2.id as "r2_id", r2.stdout as "r2_stdout", r2.stderr as "r2_stderr",
	r2.time_millis as "r2_time_millis", r2.memory_kibibytes as "r2_memory_kibibytes",
	r2.time_wall_millis as "r2_time_wall_millis", r2.exit_code as "r2_exit_code"
	from public.evaluation_test_results etr
	inner join public.runtime_data r1 on etr.exec_r_data_id = r1.id
	inner join public.runtime_data r2 on etr.checker_r_data_id = r2.id
	where etr.evaluation_id = 46
*/
