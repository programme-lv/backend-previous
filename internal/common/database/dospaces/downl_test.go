package dospaces

import (
	"github.com/programme-lv/backend/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPresignedURL(t *testing.T) {
	conf, err := config.LoadConfig(".env")
	assert.Nil(t, err)

	s3ConnParams := DOSpacesConnParams{
		AccessKey: conf.S3.Key,
		SecretKey: conf.S3.Secret,
		Region:    "fra1",
		Endpoint:  conf.S3.Endpoint,
		Bucket:    conf.S3.Bucket,
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
