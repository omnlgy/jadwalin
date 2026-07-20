package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/omnlgy/jadwalin/internal/config"
)

var MinioClient *minio.Client

func InitMinioClient(ctx context.Context, cfg *config.Config) (*minio.Client, error) {
	client, err := minio.New(cfg.MINIO_HOST, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MINIO_ACCESS_KEY, cfg.MINIO_SECRET_KEY, ""),
		Secure: cfg.MINIO_USE_SSL,
	})
	if err != nil {
		return nil, err
	}

	MinioClient = client

	exist, err := client.BucketExists(ctx, "jadwalin")
	if err != nil {
		return nil, err
	}
	if !exist {
		err = client.MakeBucket(ctx, "jadwalin", minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}
