package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// @func presignURL
// @description 서명된 다운로드 URL을 생성한다

type PresignURLInput struct {
	Bucket    string
	Key       string
	ExpiresIn int // 초 단위 (기본 3600)
	Endpoint  string
	Region    string
}

type PresignURLOutput struct {
	URL string
}

func PresignURL(in PresignURLInput) (PresignURLOutput, error) {
	client, err := newS3Client(in.Endpoint, in.Region)
	if err != nil {
		return PresignURLOutput{}, err
	}
	expiresIn := in.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	presigner := s3.NewPresignClient(client)
	req, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(in.Bucket),
		Key:    aws.String(in.Key),
	}, s3.WithPresignExpires(time.Duration(expiresIn)*time.Second))
	if err != nil {
		return PresignURLOutput{}, err
	}
	return PresignURLOutput{URL: req.URL}, nil
}
