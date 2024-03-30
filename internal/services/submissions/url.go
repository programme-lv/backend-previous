package submissions

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DownlURLGetter interface {
	GetTestDownloadURL(testSHA256 string) (string, error)
}

type S3TestURLs struct {
	presignClient *s3.PresignClient
	bucketName    string
}

var _ DownlURLGetter = &S3TestURLs{}

func (s *S3TestURLs) GetTestDownloadURL(testSHA256 string) (string, error) {
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

func NewS3TestURLs(accessKey, secretKey, region, endpoint, bucket string) (DownlURLGetter, error) {
	res := &S3TestURLs{
		presignClient: nil,
		bucketName:    bucket,
	}

	credentialProvider := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentialProvider)),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}),
		),
	)

	if err != nil {
		return nil,
			fmt.Errorf("failed to load config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	res.presignClient = s3.NewPresignClient(client)

	return res, nil
}
