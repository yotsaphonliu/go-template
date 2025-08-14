package db

import (
	"errors"

	"go-template/src/core/db/postgresql/migrations"
)

func Migrate(config *Config, dryRun bool, number int, forceMigrate bool) error {
	switch config.DBType {
	case "postgres":
		return migrations.Migrate(dryRun, number, forceMigrate)
	}
	return errors.New("unsupported database type")
}
