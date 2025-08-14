package postgresql

import (
	"context"
	"fmt"
	"os"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go-template/src/core/log"
)

type PostgresqlDB struct {
	logger log.Logger
	Config *Config

	DB *pgxpool.Pool
}

func New(config *Config, logger log.Logger) (pgdb *PostgresqlDB, err error) {
	pgdb = &PostgresqlDB{
		Config: config,
	}

	pgdb.logger = logger.WithFields(log.Fields{
		"module": "db/postgresql",
	})

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.DatabaseName,
		config.SSLMode,
	)

	// logger.Debugf("PG Connection : %s", connStr)

	var connectConf, _ = pgxpool.ParseConfig(connStr)
	connectConf.MaxConns = config.MaxOpenConns
	connectConf.MaxConnIdleTime = 5 * time.Second
	connectConf.MaxConnLifetime = 300 * time.Second
	connectConf.HealthCheckPeriod = 15 * time.Second

	connectConf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   NewDatabaseLogger(&pgdb.logger),
		LogLevel: tracelog.LogLevelTrace,
	}

	// Register Decimal Data Type to PGX Pool
	connectConf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	// Set timezone to PGX runtime
	if s := os.Getenv("TZ"); s != "" {
		connectConf.ConnConfig.RuntimeParams["timezone"] = s
	}

	pgdb.DB, err = pgxpool.NewWithConfig(context.Background(), connectConf)
	if err != nil {
		pgdb.logger.Errorf("Error connecting to postgres: %+v")
		return nil, err
	}

	return pgdb, nil
}

func (pgdb *PostgresqlDB) Close() error {
	pgdb.DB.Close()
	return nil
}
