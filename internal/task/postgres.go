package task

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/pkg/errors"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/table"
	"log"
)

type postgresTaskRepoImpl struct {
	db qrm.DB
}

func (p postgresTaskRepoImpl) GetTaskByID(taskID int64) (*Task, error) {
	stmt := postgres.SELECT(table.Tasks.AllColumns).FROM(table.Tasks).
		WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID))).LIMIT(1)

	var record model.Tasks
	err := stmt.Query(p.db, &record)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}

	return p.mapTaskTableRowToTaskDomainObject(&record)
}

func (p postgresTaskRepoImpl) GetTaskVersionByID(taskVersionID int64) (*TaskVersion, error) {
	stmt := postgres.SELECT(table.TaskVersions.AllColumns).FROM(table.TaskVersions).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID))).LIMIT(1)

	var record model.TaskVersions
	err := stmt.Query(p.db, &record)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}

	return p.mapTaskVersionTableRowToTaskVersionDomainObject(&record)
}

func (p postgresTaskRepoImpl) getStatement(statementID int64, exampleSetID int64) (*Statement, error) {
	stmt := postgres.SELECT(table.MarkdownStatements.AllColumns).FROM(table.MarkdownStatements).
		WHERE(table.MarkdownStatements.ID.EQ(postgres.Int64(statementID))).LIMIT(1)

	var record model.MarkdownStatements
	err := stmt.Query(p.db, &record)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}

	examples, err := p.getStatementExamples(exampleSetID)
	if err != nil {
		return nil, err
	}

	res := &Statement{
		ID:       record.ID,
		Story:    record.Story,
		Input:    record.Input,
		Output:   record.Output,
		Examples: examples,
		Notes:    record.Notes,
	}

	return res, nil
}

func (p postgresTaskRepoImpl) getStatementExamples(exampleSetID int64) ([]*Example, error) {
	stmt := postgres.SELECT(table.StatementExamples.AllColumns).FROM(table.StatementExamples).
		WHERE(table.StatementExamples.ExampleSetID.EQ(postgres.Int64(exampleSetID)))

	var records []*model.StatementExamples
	err := stmt.Query(p.db, &records)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}

	res := make([]*Example, 0)
	for _, r := range records {
		res = append(res, &Example{
			ID:     r.ID,
			Input:  r.Input,
			Answer: r.Answer,
		})
	}

	return res, nil
}

func (p postgresTaskRepoImpl) GetTaskVersionStatement(taskVersionID int64) (*Statement, error) {
	// TODO: implement me
	panic("implement me")
}

func (p postgresTaskRepoImpl) mapTaskVersionTableRowToTaskVersionDomainObject(taskVersionRow *model.TaskVersions) (*TaskVersion, error) {
	if taskVersionRow == nil {
		return nil, nil
	}

	var statement *Statement
	if taskVersionRow.MdStatementID != nil {
		if taskVersionRow.ExampleSetID == nil {
			return nil, errors.New("task version has statement but no example set")
		}
		var err error
		statement, err = p.getStatement(*taskVersionRow.MdStatementID, *taskVersionRow.ExampleSetID)
		if err != nil {
			return nil, err
		}
	}
	return &TaskVersion{
		ID:            taskVersionRow.ID,
		TaskID:        taskVersionRow.TaskID,
		Code:          taskVersionRow.ShortCode,
		Name:          taskVersionRow.FullName,
		Statement:     statement,
		TimeLimitMs:   taskVersionRow.TimeLimMs,
		MemoryLimitKb: taskVersionRow.MemLimKibibytes,
		CreatedAt:     taskVersionRow.CreatedAt,
	}, nil
}

func (p postgresTaskRepoImpl) mapTaskTableRowToTaskDomainObject(taskRow *model.Tasks) (*Task, error) {
	if taskRow == nil {
		return nil, nil
	}
	var current *TaskVersion
	if taskRow.CurrentVersionID != nil {
		var err error
		current, err = p.GetTaskVersionByID(*taskRow.CurrentVersionID)
		if err != nil {
			return nil, err
		}
	}
	var stable *TaskVersion
	if taskRow.StableVersionID != nil {
		var err error
		stable, err = p.GetTaskVersionByID(*taskRow.StableVersionID)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("taskRow: %+v", taskRow)
	log.Printf("current: %+v", current)
	log.Printf("stable: %+v", stable)
	return &Task{
		ID:        taskRow.ID,
		OwnerID:   taskRow.CreatedByID,
		Current:   current,
		Stable:    stable,
		CreatedAt: taskRow.CreatedAt,
	}, nil
}

func (p postgresTaskRepoImpl) ListAllTasks() ([]*Task, error) {
	stmt := postgres.SELECT(table.Tasks.AllColumns).FROM(table.Tasks)

	var records []model.Tasks
	err := stmt.Query(p.db, &records)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}

	tasks := make([]*Task, 0)
	for _, record := range records {
		task, errMappingTask := p.mapTaskTableRowToTaskDomainObject(&record)
		if errMappingTask != nil {
			return nil, errMappingTask
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (p postgresTaskRepoImpl) ListPublishedTasks() ([]*Task, error) {
	stmt := postgres.SELECT(table.Tasks.AllColumns).FROM(table.Tasks.
		INNER_JOIN(table.PublishedTaskCodes, table.PublishedTaskCodes.TaskID.EQ(table.Tasks.ID)))
	var records []model.Tasks
	err := stmt.Query(p.db, &records)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}

	tasks := make([]*Task, 0)
	for _, record := range records {
		task, errMappingTask := p.mapTaskTableRowToTaskDomainObject(&record)
		if errMappingTask != nil {
			return nil, errMappingTask
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (p postgresTaskRepoImpl) DoesTaskWithPublishedCodeExist(taskPublishedCode string) (bool, error) {
	stmt := postgres.SELECT(postgres.COUNT(table.PublishedTaskCodes.TaskID).AS("count")).
		FROM(table.PublishedTaskCodes).
		WHERE(table.PublishedTaskCodes.TaskCode.EQ(postgres.String(taskPublishedCode)))

	var record struct {
		Count int
	}
	err := stmt.Query(p.db, &record)
	if err != nil {
		return false, err
	}

	return record.Count > 0, nil
}

func (p postgresTaskRepoImpl) GetPublishedTask(taskPublishedCode string) (*Task, error) {
	stmt := postgres.SELECT(table.Tasks.AllColumns).FROM(table.Tasks.
		INNER_JOIN(table.PublishedTaskCodes, table.PublishedTaskCodes.TaskID.EQ(table.Tasks.ID))).
		WHERE(table.PublishedTaskCodes.TaskCode.EQ(postgres.String(taskPublishedCode)))

	var record model.Tasks
	err := stmt.Query(p.db, &record)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}

	return p.mapTaskTableRowToTaskDomainObject(&record)
}

func (p postgresTaskRepoImpl) MarkAsDeleted(taskID int64) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresTaskRepoImpl) UpdateStatement(taskID int64, statement *Statement) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresTaskRepoImpl) UpdateTaskNameAndCode(taskID int64, taskName string, taskCode string) error {
	//TODO implement me
	panic("implement me")
}

var _ taskRepo = postgresTaskRepoImpl{}