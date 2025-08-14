package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"go-template/src/core/log"
)

type PostgresLogger struct {
	Logger log.Logger
}

func NewDatabaseLogger(logger *log.Logger) *PostgresLogger {
	return &PostgresLogger{Logger: *logger}
}

func (pglog *PostgresLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {

	// idea from https://github.com/jackc/pgx-logrus/blob/master/adapter.go
	var logger = pglog.Logger
	if data != nil {
		logger = logger.WithFields(data)
	}

	switch level {
	case tracelog.LogLevelTrace:
		logger.WithFields(createFields("PGX_LOG_LEVEL", level)).Debugf(msg)
	case tracelog.LogLevelDebug:
		logger.Debugf(msg)
	case tracelog.LogLevelInfo:
		logger.Infof(msg)
	case tracelog.LogLevelWarn:
		logger.Warnf(msg)
	case tracelog.LogLevelError:
		logger.Errorf(msg)
	default:
		logger.WithFields(createFields("INVALID_PGX_LOG_LEVEL", level)).Errorf(msg)
	}
}

func createFields(key string, value interface{}) log.Fields {
	var fieldMap = make(map[string]interface{})
	fieldMap[key] = value
	return fieldMap
}
