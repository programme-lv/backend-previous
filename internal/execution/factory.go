package execution

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database"
)

type ExecuterFactory struct {
	DB *sqlx.DB
}

func (f *ExecuterFactory) NewExecuter(langId string, code string) (Executer, error) {
	// get programming language
	var lang database.ProgrammingLanguage
	err := f.DB.Get(&lang, "SELECT * FROM programming_languages WHERE id = $1", langId)
	if err != nil {
		return nil, err
	}

	// create temporary directory
	dir, err := os.MkdirTemp("", "programme-lv")
	if err != nil {
		return nil, err
	}

	// create code file
	file, err := os.Create(filepath.Join(dir, lang.CodeFilename))
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	// defer the closing of the file
	defer file.Close()

	// write text to file
	_, err = file.WriteString(code)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	// Save changes to file
	err = file.Sync()
	if err != nil {
		log.Fatalf("failed syncing file: %s", err)
	}

	if lang.CompileCmd != nil {
		// run compile command
		cmd := exec.Command(strings.Split(*lang.CompileCmd, " ")[0], strings.Split(*lang.CompileCmd, " ")[1:]...)
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			return nil, err
		}
	}

	return &Executable{
		directory:  dir,
		executeCmd: lang.ExecuteCmd,
	}, nil
}
