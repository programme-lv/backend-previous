package graphql

import (
	"context"
	"fmt"

	"github.com/programme-lv/backend/internal/database"
	"github.com/programme-lv/backend/internal/services/objects"
	"golang.org/x/crypto/bcrypt"
)

func (r *Resolver) GetUserFromContext(ctx context.Context) (*database.User, error) {
	userId, ok := r.SessionManager.Get(ctx, "user_id").(int64)
	if !ok {
		return nil, fmt.Errorf("user is not logged in")
	}

	var user database.User
	err := r.PostgresDB.Get(&user, "SELECT * FROM users WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Resolver) HashPassword(password string) (string, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func internalTaskVToGQLTaskV(taskVersion *objects.TaskVersion) (*TaskVersion, error) {
	if taskVersion == nil {
		return nil, nil
	}

	marshalledCreatedAt, err := taskVersion.CreatedAt.MarshalText()
	if err != nil {
		return nil, err
	}

	description, err := internalDescriptionToGQLDescription(taskVersion.Description)
	if err != nil {
		return nil, err
	}

	res := TaskVersion{
		VersionID:   fmt.Sprint(taskVersion.ID),
		Code:        taskVersion.Code,
		Name:        taskVersion.Name,
		Description: description,
		Constraints: &Constraints{
			TimeLimitMs:   int(taskVersion.TimeLimitMs),
			MemoryLimitKb: int(taskVersion.MemoryLimitKb),
		},
		Metadata:  &Metadata{},
		CreatedAt: string(marshalledCreatedAt),
	}

	return &res, nil
}

func internalDescriptionToGQLDescription(description *objects.Description) (*Description, error) {
	if description == nil {
		return nil, nil
	}

	var examples []*Example
	for _, example := range description.Examples {
		examples = append(examples, &Example{
			Input:  example.Input,
			Answer: example.Answer,
		})
	}

	res := Description{
		Story:    description.Story,
		Input:    description.Input,
		Output:   description.Output,
		Examples: examples,
		Notes:    description.Notes,
	}
	return &res, nil
}

func internalTaskToGQLTask(task *objects.Task) (*Task, error) {
	currentTaskVersion, err := internalTaskVToGQLTaskV(task.Current)
	if err != nil {
		return nil, err
	}

	stableTaskVersion, err := internalTaskVToGQLTaskV(task.Stable)
	if err != nil {
		return nil, err
	}

	marshalledCreatedAt, err := task.CreatedAt.MarshalText()
	if err != nil {
		return nil, err
	}

	marshalledUpdatedAt, err := task.UpdatedAt.MarshalText()
	if err != nil {
		return nil, err
	}

	res := Task{
		TaskID:    fmt.Sprint(task.ID),
		Current:   currentTaskVersion,
		Stable:    stableTaskVersion,
		CreatedAt: string(marshalledCreatedAt),
		UpdatedAt: string(marshalledUpdatedAt),
	}
	return &res, nil
}
