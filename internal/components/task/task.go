package task

import (
	"fmt"
	"github.com/programme-lv/backend/internal/components/user"
	"github.com/programme-lv/backend/internal/domain"
	"log/slog"
)

type Service interface {
	// ListPublishedTasks returns a list of all tasks that have been published.
	// Not all tasks have been published, furthermore there may exist unpublished
	// tasks with the same name and code as published ones. A task can be published
	// by an administrator. Usually it is quality checked and reviewed before being.
	ListPublishedTasks() ([]*domain.Task, error)

	// GetTaskByID returns the task with the given ID with both the
	// current and the stable version populated. Method is accessible only
	// to the task creator or an administrator since it may contain sensitive information.
	// For retrieving user facing information use GetPublicTaskVersionByPublishedCode.
	GetTaskByID(actingUserID, taskID int64) (*domain.Task, error)

	// GetPublicTaskVersionByPublishedCode finds the task with the published code
	// and returns its "STABLE" version. This method is used to retrieve user facing
	// information about a task. It is accessible to all users.
	GetPublicTaskVersionByPublishedCode(taskPublishedCode string) (*domain.TaskVersion, error)

	// ListEditableTasks returns a list of all tasks that the user can edit.
	// An administrator can edit all tasks, while others can only edit tasks they have created.
	// In the future task edit permission sharing may be implemented.
	ListEditableTasks(actingUserID int64) ([]*domain.Task, error)

	// ListUserSolvedTasks returns a list of all tasks that the user has solved.
	// As of now the user is considered to have solved the task if they have submitted
	// a solution that has maximum score.
	ListUserSolvedTasks(actingUserID int64) ([]*domain.Task, error)

	// CreateTask creates a new task with the given name and code.
	// All unprovided fields are set to their default values such as an example statement.
	CreateTask(actingUserID int64, taskCode string, taskName string) (*domain.Task, error)

	// UpdateTaskStatement finds the "CURRENT" task version, duplicates it
	// (creating a new task version) and updates the statement with the new one.
	UpdateTaskStatement(actingUserID int64, taskID int64, statement *domain.Statement) error

	// UpdateTaskNameAndCode finds the "CURRENT" task version, duplicates it
	// (creating a new task version) and updates the name and code with the new ones.
	UpdateTaskNameAndCode(actingUserID int64, taskID int64, taskName string, taskCode string) error

	// DeleteTask updates the DeletedAt of the task to the current time
	// hiding it from future queries.
	DeleteTask(actingUserID int64, taskID int64) error
}

type taskRepo interface {
	GetUserByID(userID int64) (*domain.User, error)

	ListPublishedTasks() ([]*domain.Task, error)
	GetTaskByID(taskID int64) (*domain.Task, error)
	GetPublishedTask(taskPublishedCode string) (*domain.Task, error)
	MarkAsDeleted(taskID int64) error

	// UpdateStatement duplicates the current task version, creates a new statement
	// and updates the task version with the new statement. All in one transaction.
	UpdateStatement(taskID int64, statement *domain.Statement) error

	// UpdateTaskNameAndCode duplicates the current task version, creates a new task version
	// and updates the task version with the new name and code. All in one transaction.
	UpdateTaskNameAndCode(taskID int64, taskName string, taskCode string) error
}

type service struct {
	logger  *slog.Logger
	userSrv user.Service
	repo    taskRepo
}

func (s service) ListPublishedTasks() ([]*domain.Task, error) {
	return s.repo.ListPublishedTasks()
}

func (s service) GetTaskByID(actingUserID, taskID int64) (*domain.Task, error) {
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
		return nil, err

	}
	return s.repo.GetTaskByID(taskID)
}

func (s service) GetPublicTaskVersionByPublishedCode(taskPublishedCode string) (*domain.TaskVersion, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) ListEditableTasks(actingUserID int64) ([]*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) ListUserSolvedTasks(actingUserID int64) ([]*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) CreateTask(actingUserID int64, taskCode string, taskName string) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdateTaskStatement(actingUserID int64, taskID int64, statement *domain.Statement) error {
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
