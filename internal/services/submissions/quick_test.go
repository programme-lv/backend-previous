package submissions_test

import (
	"testing"

	"github.com/programme-lv/backend/config"
	"github.com/programme-lv/backend/internal/services/submissions"
)

func TestGetEvaluationObj(t *testing.T) {
	db, err := config.ConnectToPostgresByEnvConf()
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	evalObj, err := submissions.GetEvaluationObj(db, 45, true)
	if err != nil {
		t.Fatalf("failed to get evaluation obj: %v", err)
	}

	t.Logf("evaluation obj: %+v", evalObj)
}
