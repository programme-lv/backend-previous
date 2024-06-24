package dospaces

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DOSpacesS3ObjStorage struct {
	presignClient *s3.PresignClient
	bucketName    string
}

//var _ eval.TestDownloadURLProvider = &DOSpacesS3ObjStorage{}
