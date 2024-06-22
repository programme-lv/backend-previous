package graphql

import (
	"fmt"
	"github.com/programme-lv/backend/internal/lang"
	"github.com/programme-lv/backend/internal/task"
	"github.com/programme-lv/backend/internal/user"
)

func mapDomainUserObjToGQLUserObj(user *user.User) *User {
	return &User{
		ID:        fmt.Sprint(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
	}
}

func internalTaskVToGQLTaskV(taskVersion *task.TaskVersion) (*TaskVersion, error) {
	if taskVersion == nil {
		return nil, nil
	}

	marshalledCreatedAt, err := taskVersion.CreatedAt.MarshalText()
	if err != nil {
		return nil, err
	}

	description, err := internalDescriptionToGQLDescription(taskVersion.Statement)
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
		CreatedAt: string(marshalledCreatedAt),
	}

	return &res, nil
}

func internalDescriptionToGQLDescription(description *task.Statement) (*Description, error) {
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

func mapDomainTaskObjToGQLTask(task *task.Task) (*Task, error) {
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

	res := Task{
		TaskID:    fmt.Sprint(task.ID),
		Current:   currentTaskVersion,
		Stable:    stableTaskVersion,
		CreatedAt: string(marshalledCreatedAt),
	}
	return &res, nil
}

func internalProgrammingLanguageToGraphQL(lang *lang.ProgrammingLanguage) *ProgrammingLanguage {
	return &ProgrammingLanguage{
		ID:       lang.ID,
		FullName: lang.Name,
		MonacoID: lang.MonacoID,
		Enabled:  true,
	}
}
