package dospaces

import (
	"log/slog"
	"testing"

	"github.com/programme-lv/backend/internal/environment"
)

func TestGetPresignedURL(t *testing.T) {
	config := environment.ReadEnvConfig(slog.Default())

	s3ConnParams := DOSpacesConnParams{
		AccessKey: config.DOSpacesKey,
		SecretKey: config.DOSpacesSecret,
		Region:    "fra1",
		Endpoint:  config.S3Endpoint,
		Bucket:    config.S3Bucket,
	}
	urls, err := NewDOSpacesConn(s3ConnParams)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	res, err := urls.GetTestDownloadURL("241aebeb0c7f29cae343e0da0366b2ee92a026067dd4484358659fa4e0cd84a4")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Log(res)
}
