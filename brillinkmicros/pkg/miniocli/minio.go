package miniocli

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

type McInterface interface {
	GetCli() *minio.Client
	DefaultMakeBucketOptions() minio.MakeBucketOptions
	UploadFile(ctx context.Context, bucketName string, objName string, reader io.Reader) error
	DownloadFile(ctx context.Context, bucketName string, objName string) (io.ReadCloser, error)
}

type MinioClient struct {
	Cli *minio.Client
}

func (mc MinioClient) DefaultMakeBucketOptions() minio.MakeBucketOptions {
	return minio.MakeBucketOptions{
		Region: "ASIA-EAST2",
	}
}

func (mc MinioClient) UploadFile(ctx context.Context, bucketName string, objName string, reader io.Reader) error {
	exists, errBucketExists := mc.Cli.BucketExists(ctx, bucketName)
	if errBucketExists != nil {
		return errBucketExists
	}
	if !exists {
		err := mc.Cli.MakeBucket(ctx, bucketName, mc.DefaultMakeBucketOptions())
		if err != nil {
			return err
		}
	}
	_, err := mc.Cli.PutObject(ctx, bucketName, objName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (mc MinioClient) GetCli() *minio.Client {
	return mc.Cli
}
