package minio

import (
	"context"
	"go-template/src/core/log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO interface {
	CreateObject(ctx context.Context, objectName string, body []byte) error
	DownloadFile(ctx context.Context, objectName string) ([]byte, error)
	CreateDefaultBucket() error
}

type minIO struct {
	conf   *Config
	client *minio.Client
	log    log.Logger
}

// New creates a new minIO client using the provided configuration.
func New(logger log.Logger) (MinIO, error) {
	conf, err := InitConfig()
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(conf.EndpointUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Username, conf.Password, ""),
		Secure: false, // Set to true if you use SSL
	})

	if err != nil {
		return nil, err
	}
	return &minIO{
		conf:   conf,
		client: minioClient,
		log:    logger,
	}, nil
}

func (m *minIO) defaultBucket() string {
	return m.conf.Bucket
}
