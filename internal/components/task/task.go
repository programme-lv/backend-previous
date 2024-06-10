package task

import "github.com/programme-lv/backend/internal/domain"

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
	GetTaskByID(taskID int64) (*domain.Task, error)

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
