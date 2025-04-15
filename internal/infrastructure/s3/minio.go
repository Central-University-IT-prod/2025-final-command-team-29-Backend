package s3

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"backend/internal/infrastructure"
	"backend/internal/infrastructure/logger"
)

func NewMinioClient(config *infrastructure.Config, log *logger.Logger) (*minio.Client, string, error) {
	minioClient, err := minio.New(config.S3.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3.S3Assets, config.S3.S3Secret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %s", err.Error())
		return nil, "", err
	}

	log.Infof("Connected to MinIO at %s", config.S3.S3Endpoint)
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, config.S3.S3Backet)
	if err != nil {
		log.Fatalf("Failed to check if bucket exists: %s", err.Error())
		return nil, "", err
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, config.S3.S3Backet, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create bucket: %s", err.Error())
			return nil, "", err
		}
		log.Infof("Bucket %s created successfully", config.S3.S3Backet)
	} else {
		log.Infof("Bucket %s already exists", config.S3.S3Backet)
	}

	return minioClient, config.S3.S3Backet, nil
}
