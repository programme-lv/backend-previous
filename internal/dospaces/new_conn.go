package dospaces

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DOSpacesConnParams struct {
	AccessKey string
	SecretKey string
	Region    string
	Endpoint  string
	Bucket    string
}

func NewDOSpacesConn(params DOSpacesConnParams) (*DOSpacesS3ObjStorage, error) {
	accessKey := params.AccessKey
	secretKey := params.SecretKey
	region := params.Region
	endpoint := params.Endpoint
	bucket := params.Bucket

	res := &DOSpacesS3ObjStorage{
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
