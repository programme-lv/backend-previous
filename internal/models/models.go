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
