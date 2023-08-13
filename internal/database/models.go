package database

import "time"

type User struct {
	ID             int64      `db:"id"`
	Username       string     `db:"username"`
	HashedPassword string     `db:"hashed_password"`
	Email          string     `db:"email"`
	FirstName      string     `db:"first_name"`
	LastName       string     `db:"last_name"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
	IsAdmin        bool       `db:"is_admin"`
}

type ProgrammingLanguage struct {
	ID             string  `db:"id"`
	FullName       string  `db:"full_name"`
	CodeFilename   string  `db:"code_filename"`
	CompileCmd     *string `db:"compile_cmd"`
	ExecuteCmd     string  `db:"execute_cmd"`
	EnvVersionCmd  string  `db:"env_version_cmd"`
	HelloWorldCode string  `db:"hello_world_code"`
	MonacoId       string  `db:"monaco_id"`
}

type Task struct {
	ID          int64     `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	CreatedByID int64     `db:"created_by_id"`

	RelevantVersionID  *int64 `db:"relevant_version_id"`
	PublishedVersionID *int64 `db:"published_version_id"`
}

type TaskVersion struct {
	ID                int64      `db:"id"`
	TaskID            int64      `db:"task_id"`
	ShortCode         string     `db:"short_code"`
	FullName          string     `db:"full_name"`
	TimeLimMs         int        `db:"time_lim_ms"`
	MemLimKb          int        `db:"mem_lim_kb"`
	TestingTypeID     string     `db:"testing_type_id"`
	Origin            *string    `db:"origin"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
	CheckerTextID     *int64     `db:"checker_text_id"`
	InteratctorTextID *int64     `db:"interactor_text_id"`
}

type TaskAuthor struct {
	TaskID string `db:"task_id"`
	Author string `db:"author"`
}

type EvalType struct {
	ID            string `db:"id"`
	DescriptionEn string `db:"description_en"`
}

type TaskSource struct {
	Abbreviation string `db:"abbreviation"`
	FullName     string `db:"full_name"`
}

type MarkdownStatement struct {
	ID            int64  `db:"id"`
	Story         string `db:"story"`
	Input         string `db:"input"`
	Output        string `db:"output"`
	Notes         *string `db:"notes"`
	Scoring       *string `db:"scoring"`
	TaskVersionID int64  `db:"task_version_id"`
}
