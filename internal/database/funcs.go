package database

import "github.com/jmoiron/sqlx"

func SelectUserByUsername(db *sqlx.DB, username string) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	return &user, err
}

func DoesUserExistByUsername(db *sqlx.DB, username string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE username = $1", username)
	return count > 0, err
}

func DoesUserExistByEmail(db *sqlx.DB, email string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", email)
	return count > 0, err
}

func CreateUser(db *sqlx.DB, username string, hashed_password string, email string,
	firstName string, lastName string) error {
	_, err := db.Exec("INSERT INTO users (username, hashed_password, email, first_name, last_name, created_at) VALUES ($1, $2, $3, $4, $5, now())", username, hashed_password, email, firstName, lastName)
	return err
}

func DeleteUserById(db *sqlx.DB, id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// creates a task submission and returns its id
func CreateTaskSubmission(db sqlx.Ext, userId int64, taskId int64, programmingLangId string, submission string) (int64, error) {
	var id int64
	err := db.QueryRowx("INSERT INTO task_submissions (user_id, task_id, programming_lang_id, submission, created_at) VALUES ($1, $2, $3, $4, now()) RETURNING id", userId, taskId, programmingLangId, submission).Scan(&id)
	return id, err
}

func GetTaskById(db *sqlx.DB, id int64) (*Task, error) {
	var task Task
	err := db.Get(&task, "SELECT * FROM tasks WHERE id = $1", id)
	return &task, err
}

func CreateSubmissionEvaluation(db sqlx.Execer,
	taskSubmissionId int64, evalTaskVersionId int64,
	testMaximumTimeMs *int64, testMaximumMemoryKb *int64,
	testTotalTimeMs int64, testTotalMemoryKb int64,
	evalStatusId string, evalTotalScore int64, evalPossibleScore int64,
	compilationStdout *string, compilationStderr *string,
	compilationTimeMs *int64, compilationMemoryKb *int64) error {
	_, err := db.Exec(`INSERT INTO submission_evaluations
	(task_submission_id, eval_task_version_id,
	test_maximum_time_ms, test_maximum_memory_kb,
	test_total_time_ms, test_total_memory_kb,
	eval_status_id, eval_total_score, eval_possible_score,
	compilation_stdout, compilation_stderr,
	compilation_time_ms, compilation_memory_kb,
	created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, now())`,
		taskSubmissionId, evalTaskVersionId,
		testMaximumTimeMs, testMaximumMemoryKb,
		testTotalTimeMs, testTotalMemoryKb, evalStatusId,
		evalTotalScore, evalPossibleScore,
		compilationStdout, compilationStderr,
		compilationTimeMs, compilationMemoryKb)
	return err
}
