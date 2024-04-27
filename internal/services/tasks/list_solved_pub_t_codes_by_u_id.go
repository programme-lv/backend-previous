package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func ListSolvedPublishedTaskCodesByUserID(db qrm.DB, userID int64) ([]string, error) {
	/*
		select distinct task_code from published_task_codes ptc
		inner join task_submissions ts on ts.task_id = ptc.task_id
		inner join evaluations e on e.id = ts.visible_eval_id
		where eval_total_score = eval_possible_score and ts.user_id = 1
	*/
	selectStmt := postgres.SELECT(table.PublishedTaskCodes.TaskCode).DISTINCT().
		FROM(table.PublishedTaskCodes.
			INNER_JOIN(table.TaskSubmissions, table.TaskSubmissions.TaskID.EQ(table.PublishedTaskCodes.TaskID)).
			INNER_JOIN(table.Evaluations, table.Evaluations.ID.EQ(table.TaskSubmissions.VisibleEvalID))).
		WHERE(table.Evaluations.EvalTotalScore.EQ(table.Evaluations.EvalPossibleScore).
			AND(table.TaskSubmissions.UserID.EQ(postgres.Int64(userID))))
	var records []model.PublishedTaskCodes
	err := selectStmt.Query(db, &records)
	if err != nil {
		return nil, err
	}

	taskCodes := make([]string, len(records))
	for i, record := range records {
		taskCodes[i] = record.TaskCode
	}
	return taskCodes, nil
}
