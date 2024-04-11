package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func ListSolvedTasksByUserID(db qrm.Queryable, userID int64) ([]int64, error) {
	stmt := postgres.SELECT(table.TaskSubmissions.TaskID).
		FROM(table.TaskSubmissions.INNER_JOIN(table.Evaluations,
			table.Evaluations.ID.EQ(table.TaskSubmissions.VisibleEvalID))).
		WHERE(table.TaskSubmissions.UserID.EQ(postgres.Int64(userID)).
			AND(table.TaskSubmissions.Hidden.EQ(postgres.Bool(false))).
			AND(table.Evaluations.EvalPossibleScore.IS_NOT_NULL()).
			AND(table.Evaluations.EvalTotalScore.EQ(table.Evaluations.EvalPossibleScore))).
		GROUP_BY(table.TaskSubmissions.TaskID)

	var taskIDs []int64
	err := stmt.Query(db, &taskIDs)
	if err != nil {
		return nil, err
	}

	return taskIDs, nil
}
