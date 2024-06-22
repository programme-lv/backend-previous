package dospaces

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/programme-lv/backend/internal/eval"
)

type DOSpacesS3ObjStorage struct {
	presignClient *s3.PresignClient
	bucketName    string
}

var _ eval.TestDownloadURLProvider = &DOSpacesS3ObjStorage{}
