package models

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
	ID        string    `db:"id"`
	FullName  string    `db:"full_name"`
	Origin    *string   `db:"origin"`
	CreatedAt time.Time `db:"created_at"`
}

type TaskVersion struct {
	ID          int64      `db:"id"`
	VersionName string     `db:"version_name"`
	TaskID      string     `db:"task_id"`
	TimeLimMs   int        `db:"time_lim_ms"`
	MemLimKb    int        `db:"mem_lim_kb"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	EvalTypeID  string     `db:"eval_type_id"`
}

type TaskAuthor struct {
	TaskID string `db:"task_id"`
	Author string `db:"author"`
}

type EvalType struct {
	ID            string `db:"id"`
	DescriptionEn string `db:"description_en"`
}
