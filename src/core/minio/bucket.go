package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
)

func (m *minIO) CreateDefaultBucket() error {
	err := m.client.MakeBucket(
		context.Background(),
		m.defaultBucket(),
		minio.MakeBucketOptions{Region: "thailand"},
	)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := m.client.BucketExists(context.Background(), m.defaultBucket())
		if errBucketExists == nil && exists {
			m.log.Infof("Bucket: %s already exists", m.defaultBucket())
			return nil
		}

		return err
	}

	m.log.Infof("Successfully created %s object storage", m.defaultBucket())

	return nil
}
