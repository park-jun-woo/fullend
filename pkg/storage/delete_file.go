package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// @func deleteFile
// @description S3 호환 스토리지에서 파일을 삭제한다

type DeleteFileInput struct {
	Bucket   string
	Key      string
	Endpoint string
	Region   string
}

type DeleteFileOutput struct{}

func DeleteFile(in DeleteFileInput) (DeleteFileOutput, error) {
	client, err := newS3Client(in.Endpoint, in.Region)
	if err != nil {
		return DeleteFileOutput{}, err
	}
	_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(in.Bucket),
		Key:    aws.String(in.Key),
	})
	return DeleteFileOutput{}, err
}
