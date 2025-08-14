package db

import (
	"time"

	"go-template/src/core/model"
)

type DBApiKeysInterface interface {
	InsertApiKeys(apiKeysList []model.ApiKey, isRoot bool) error
	VerifyApiKey(key string, newExpireTime time.Time) ([]*model.ApiKey, error)
	DeleteExpireApiKey() error
	DeleteApiKey(key string) error
}
