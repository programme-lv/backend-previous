package graphql

import (
	"fmt"
	"github.com/programme-lv/backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

func (r *Resolver) HashPassword(password string) (string, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func internalTaskVToGQLTaskV(taskVersion *domain.TaskVersion) (*TaskVersion, error) {
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
		Metadata:  &Metadata{},
		CreatedAt: string(marshalledCreatedAt),
	}

	return &res, nil
}

func internalDescriptionToGQLDescription(description *domain.Statement) (*Description, error) {
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

func internalTaskToGQLTask(task *domain.Task) (*Task, error) {
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

func internalSubmissionToGQLSubmission(submission *domain.TaskSubmission) (*Submission, error) {
	marshalledCreatedAt, err := submission.CreatedAt.MarshalText()
	if err != nil {
		return nil, err
	}

	res := Submission{
		ID:         fmt.Sprint(submission.ID),
		Task:       nil,
		Language:   nil,
		Submission: submission.Content,
		Evaluation: nil,
		Username:   submission.Author.Username,
		CreatedAt:  string(marshalledCreatedAt),
	}

	GQLTask, err := internalTaskToGQLTask(submission.Task)
	if err != nil {
		return nil, err
	}
	res.Task = GQLTask

	GQLLang := internalProgrammingLanguageToGraphQL(submission.Language)
	res.Language = GQLLang

	if submission.VisibleEval != nil {
		visEvalObj, err := internalEvalObjToGQLEvaluation(submission.VisibleEval)
		if err != nil {
			return nil, err
		}
		res.Evaluation = visEvalObj
	}

	return &res, nil
}

func internalEvalObjToGQLEvaluation(eval *domain.Evaluation) (*Evaluation, error) {
	var possibleScoreInt32 *int
	if eval.PossibleScore != nil {
		possibleScoreInt32 = new(int)
		*possibleScoreInt32 = int(*eval.PossibleScore)
	}

	res := Evaluation{
		ID:                fmt.Sprint(eval.ID),
		Status:            eval.StatusID,
		TotalScore:        int(eval.ReceivedScore),
		PossibleScore:     possibleScoreInt32,
		RuntimeStatistics: nil, // TODO
		CompileRData:      internalRDataToGQLRData(eval.CheckerRunData),
		TestResults:       nil,
	}

	var testResults []*TestResult = make([]*TestResult, len(eval.TestResults))
	for i, testResult := range eval.TestResults {
		testResults[i] = &TestResult{
			ID:            fmt.Sprint(testResult.ID),
			TaskVTestID:   fmt.Sprint(testResult.TaskVTestID),
			UserSubmRData: internalRDataToGQLRData(testResult.ExecRData),
			CheckerRData:  internalRDataToGQLRData(testResult.CheckerRData),
			Result:        TestResultType(testResult.EvalStatusID),
		}
	}
	res.TestResults = testResults

	return &res, nil
}

func internalRDataToGQLRData(data *domain.RuntimeData) *RuntimeData {
	if data == nil {
		return nil
	}
	res := &RuntimeData{
		TimeMs:   0,
		MemoryKb: 0,
		ExitCode: 0,
		Stdout:   "",
		Stderr:   "",
	}

	if data.ExitCode != nil {
		res.ExitCode = int(*data.ExitCode)
	}
	if data.MemoryKibibytes != nil {
		res.MemoryKb = int(*data.MemoryKibibytes)
	}
	if data.TimeMillis != nil {
		res.TimeMs = int(*data.TimeMillis)
	}
	if data.Stdout != nil {
		res.Stdout = *data.Stdout
	}
	if data.Stderr != nil {
		res.Stderr = *data.Stderr
	}

	return res
}

func internalProgrammingLanguageToGraphQL(lang *domain.ProgrammingLanguage) *ProgrammingLanguage {
	return &ProgrammingLanguage{
		ID:       lang.ID,
		FullName: lang.Name,
		MonacoID: lang.MonacoID,
		Enabled:  true,
	}
}
