package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// @func uploadFile
// @description S3 호환 스토리지에 파일을 업로드한다

type UploadFileInput struct {
	Bucket      string
	Key         string
	Data        []byte
	ContentType string
	Endpoint    string // MinIO 등 커스텀 엔드포인트 (빈 문자열이면 AWS 기본)
	Region      string
}

type UploadFileOutput struct {
	URL string
}

func UploadFile(in UploadFileInput) (UploadFileOutput, error) {
	client, err := newS3Client(in.Endpoint, in.Region)
	if err != nil {
		return UploadFileOutput{}, err
	}
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(in.Bucket),
		Key:         aws.String(in.Key),
		Body:        bytes.NewReader(in.Data),
		ContentType: aws.String(in.ContentType),
	})
	if err != nil {
		return UploadFileOutput{}, err
	}
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", in.Bucket, in.Region, in.Key)
	if in.Endpoint != "" {
		url = fmt.Sprintf("%s/%s/%s", in.Endpoint, in.Bucket, in.Key)
	}
	return UploadFileOutput{URL: url}, nil
}
