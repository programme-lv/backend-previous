package task

import (
	"fmt"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/user"
	"log/slog"
)

type Service interface {
	// ListPublishedTasks returns a list of all tasks that have been published.
	// Not all tasks have been published, furthermore there may exist unpublished
	// tasks with the same name and code as published ones. A task can be published
	// by an administrator. Usually it is quality checked and reviewed before being.
	ListPublishedTasks() ([]*Task, error)

	// GetTaskByID returns the task with the given ID with both the
	// current and the stable version populated. Method is accessible only
	// to the task creator or an administrator since it may contain sensitive information.
	GetTaskByID(actingUserID, taskID int64) (*Task, error)

	// GetTaskByPublishedCode finds the task with the published code
	// This method is used to retrieve user facing information about a task.
	GetTaskByPublishedCode(taskPublishedCode string) (*Task, error)

	// ListEditableTasks returns a list of all tasks that the user can edit.
	// An administrator can edit all tasks, while others can only edit tasks they have created.
	// In the future task edit permission sharing may be implemented.
	ListEditableTasks(actingUserID int64) ([]*Task, error)

	// CreateTask creates a new task with the given name and code.
	// All unprovided fields are set to their default values such as an example statement.
	CreateTask(actingUserID int64, taskCode string, taskName string) (*Task, error)

	// UpdateTaskStatement finds the "CURRENT" task version, duplicates it
	// (creating a new task version) and updates the statement with the new one.
	UpdateTaskStatement(actingUserID int64, taskID int64, statement *Statement) error

	// UpdateTaskNameAndCode finds the "CURRENT" task version, duplicates it
	// (creating a new task version) and updates the name and code with the new ones.
	UpdateTaskNameAndCode(actingUserID int64, taskID int64, taskName string, taskCode string) error

	// DeleteTask updates the DeletedAt of the task to the current time
	// hiding it from future queries.
	DeleteTask(actingUserID int64, taskID int64) error
}

type taskRepo interface {
	GetTaskByID(taskID int64) (*Task, error)

	ListAllTasks() ([]*Task, error)
	ListPublishedTasks() ([]*Task, error)

	DoesTaskWithPublishedCodeExist(taskPublishedCode string) (bool, error)
	GetPublishedTask(taskPublishedCode string) (*Task, error)

	MarkAsDeleted(taskID int64) error

	// UpdateStatement duplicates the current task version, creates a new statement
	// and updates the task version with the new statement. All in one transaction.
	UpdateStatement(taskID int64, statement *Statement) error

	// UpdateTaskNameAndCode duplicates the current task version, creates a new task version
	// and updates the task version with the new name and code. All in one transaction.
	UpdateTaskNameAndCode(taskID int64, taskName string, taskCode string) error
}

type service struct {
	logger  *slog.Logger
	userSrv user.Service
	repo    taskRepo
}

func NewService(userSrv user.Service, db qrm.DB) Service {
	return service{
		logger:  slog.Default().With("service", "task"),
		userSrv: userSrv,
		repo:    postgresTaskRepoImpl{db: db},
	}
}

func (s service) ListPublishedTasks() ([]*Task, error) {
	return s.repo.ListPublishedTasks()
}

func (s service) GetTaskByID(actingUserID, taskID int64) (*Task, error) {
	actingUser, err := s.userSrv.GetUserByID(actingUserID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("getting user by ID: %v", err))
		return nil, err
	}
	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("getting task by ID: %v", err))
		return nil, err
	}
	if task.OwnerID != actingUser.ID && !actingUser.IsAdmin {
		return nil, newErrorUserDoesNotHaveEditAccessToTask()

	}
	return s.repo.GetTaskByID(taskID)
}

func (s service) GetTaskByPublishedCode(taskPublishedCode string) (*Task, error) {
	exists, err := s.repo.DoesTaskWithPublishedCodeExist(taskPublishedCode)
	if err != nil {
		s.logger.Error(fmt.Sprintf("checking if task with published code exists: %v", err))
		return nil, err
	}
	if !exists {
		return nil, newErrorTaskNotFound()
	}

	task, err := s.repo.GetPublishedTask(taskPublishedCode)
	if err != nil {
		s.logger.Error(fmt.Sprintf("getting published task: %v", err))
		return nil, err
	}

	return task, nil
}

func (s service) ListEditableTasks(actingUserID int64) ([]*Task, error) {
	tasks, err := s.repo.ListAllTasks()
	if err != nil {
		s.logger.Error(fmt.Sprintf("listing all tasks: %v", err))
		return nil, err
	}

	actingUser, err := s.userSrv.GetUserByID(actingUserID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("getting user by ID: %v", err))
		return nil, err
	}

	editableTasks := make([]*Task, 0)
	for _, task := range tasks {
		if task.OwnerID == actingUser.ID || actingUser.IsAdmin {
			editableTasks = append(editableTasks, task)
		}
	}

	return editableTasks, nil
}

func (s service) CreateTask(actingUserID int64, taskCode string, taskName string) (*Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdateTaskStatement(actingUserID int64, taskID int64, statement *Statement) error {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdateTaskNameAndCode(actingUserID int64, taskID int64, taskName string, taskCode string) error {
	//TODO implement me
	panic("implement me")
}

func (s service) DeleteTask(actingUserID int64, taskID int64) error {
	//TODO implement me
	panic("implement me")
}

var _ Service = service{}
