package minio_config

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MinioConfig struct {
	Endpoint    string `env:"MINIO_ENDPOINT"`
	AccessKeyID string `env:"MINIO_ACCESS_KEY"`
	SecretKey   string `env:"MINIO_SECRET_KEY"`
	UseSSL      bool   `env:"MINIO_USE_SSL"`
	BucketName  string `env:"MINIO_BUCKET_NAME"`
}

func InitBuckets(ctx context.Context, client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func UploadFile(ctx context.Context, client *minio.Client, fileHeader *multipart.FileHeader, bucket string) (string, error) {

	minioID := uuid.New().String()

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if fileHeader.Size > 10*1024*1024 {
		return "", fmt.Errorf("file too large: %d bytes", fileHeader.Size)
	}

	_, err = client.PutObject(ctx, bucket, minioID, file, fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
	if err != nil {
		return "", fmt.Errorf("minio upload failed: %w", err)
	}

	return minioID, nil
}
