package minio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
)

func (m *minIO) CreateObject(ctx context.Context, objectName string, body []byte) error {

	reader := bytes.NewReader(body)
	contentType := http.DetectContentType(body)

	_, err := m.client.PutObject(ctx, m.defaultBucket(), objectName, reader, reader.Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		m.log.Errorf("Failed to upload file to bucket '%s': %+v", m.defaultBucket(), err)
		return err
	}

	m.log.Debugf("Object '%s' created successfully in bucket '%s'", objectName, m.defaultBucket())
	return nil
}

// DownloadFile downloads a file from the specified bucket.
func (m *minIO) DownloadFile(ctx context.Context, objectName string) ([]byte, error) {

	obj, err := m.client.GetObject(ctx, m.defaultBucket(), objectName, minio.GetObjectOptions{})
	if err != nil {
		m.log.Infof("Failed to download file: %+v", err)
		return nil, err
	}

	b, err := io.ReadAll(obj)
	if err != nil {
		m.log.Infof("Failed to read object: %+v", err)
		return nil, err
	}
	defer obj.Close()

	m.log.Debugf("File downloaded successfully: %+v", objectName)
	return b, nil
}
