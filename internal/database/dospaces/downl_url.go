package dospaces

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *DOSpacesS3ObjStorage) GetTestDownloadURL(testSHA256 string) (string, error) {
	objectKey := fmt.Sprintf("tests/%s", testSHA256)
	request, err := s.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(24 * time.Hour) // 24 hours
	})
	if err != nil {
		return "",
			fmt.Errorf("failed to presign object: %v", err)
	}
	return request.URL, nil
}
