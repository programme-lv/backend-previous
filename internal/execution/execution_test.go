package execution

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/programme-lv/backend/internal/environment"
	"github.com/programme-lv/backend/internal/database"
)

func TestProgrammingLanguages(t *testing.T) {
	conf := environment.ReadEnvConfig()
	sqlxDb := sqlx.MustConnect("postgres", conf.SqlxConnString)
	defer sqlxDb.Close()

	// fetch all languages
	var langs []database.ProgrammingLanguage
	err := sqlxDb.Select(&langs, "SELECT * FROM programming_languages")
	if err != nil {
		t.Fatal(err)
	}

	for _, lang := range langs {
		t.Run(lang.FullName, func(t *testing.T) {
			factory := ExecuterFactory{DB: sqlxDb}
			executable, err := factory.NewExecuter(lang.ID, lang.HelloWorldCode)
			if err != nil {
				t.Fatal(err)
			}
			defer executable.Cleanup()

			// execute code
			result, err := executable.Execute()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(result)

			// log lang env version cmd output
			cmd := exec.Command(strings.Split(lang.EnvVersionCmd, " ")[0], strings.Split(lang.EnvVersionCmd, " ")[1:]...)
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(stdoutStderr))
		})
	}
}

func TestExecuterFactory_NewExecuter(t *testing.T) {
	conf := environment.ReadEnvConfig()
	sqlxDb := sqlx.MustConnect("postgres", conf.SqlxConnString)
	defer sqlxDb.Close()

	factory := ExecuterFactory{DB: sqlxDb}
	executable, err := factory.NewExecuter("python3.10", "print('Hello, World!')")
	if err != nil {
		t.Fatal(err)
	}
	defer executable.Cleanup()

	// execute code
	result, err := executable.Execute()
	if err != nil {
		t.Fatal(err)
	}

	// compare result stdout to "Hello, World!\n"
	if result.Stdout != "Hello, World!\n" {
		t.Fatalf("expected stdout to be \"Hello, World!\\n\", got \"%s\"", result.Stdout)
	}

	t.Log(result)
}
