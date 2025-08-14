package db

import (
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-template/src/core/db/postgresql"
	"go-template/src/core/log"
)

type DB interface {
	DBApiKeysInterface
	DBActivityLogInterface

	Close() error
}

type PostgresqlDB struct {
	logger log.Logger
	DB     *pgxpool.Pool
	Tx     pgx.Tx
}

func New(config *Config, logger log.Logger) (db DB, err error) {
	switch config.DBType {
	case "postgres":
		dbConfig, err := postgresql.InitConfig()
		if err != nil {
			return nil, err
		}

		return postgresql.New(dbConfig, logger)
	}

	return nil, errors.New("unsupported database type")
}
