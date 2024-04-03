package submissions

import (
	"log/slog"
	"testing"

	"github.com/programme-lv/backend/internal/environment"
)

func TestGetPresignedURL(t *testing.T) {
	config := environment.ReadEnvConfig(slog.Default())

	urls, err := NewS3TestURLs(config.DOSpacesKey, config.DOSpacesSecret, "fra1", config.S3Endpoint, config.S3Bucket)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	res, err := urls.GetTestDownloadURL("241aebeb0c7f29cae343e0da0366b2ee92a026067dd4484358659fa4e0cd84a4")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Log(res)
}
