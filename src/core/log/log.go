package log

import (
	"errors"
)

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	//Debug has verbose message
	Debug = "debug"
	//Info is default log level
	Info = "info"
	//Warn is for logging messages about possible issues
	Warn = "warn"
	//Error is for logging errors
	Error = "error"
	//Fatal is for logging fatal messages. The system shutsdown after logging the message.
	Fatal = "fatal"
)

const (
	InstanceZapLogger int = iota
	InstanceLogrusLogger
)

var (
	errInvalidLoggerInstance = errors.New("invalid logger instance")
)

// Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

// NewLogger returns an instance of logger
func NewLogger(config *Configuration, loggerInstance int) (Logger, error) {
	if config == nil {
		config = &Configuration{
			EnableConsole:     true,
			ConsoleLevel:      "debug",
			ConsoleJSONFormat: false,
			EnableFile:        false,
			// FileLevel:         log.Info,
			// FileJSONFormat:    true,
			// FileLocation:      "log.log",
			Color: true,
		}
	}

	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(*config)
		if err != nil {
			return nil, err
		}
		return logger, nil
	case InstanceLogrusLogger:
		logger, err := newLogrusLogger(*config)
		if err != nil {
			return nil, err
		}
		return logger, nil

	default:
		return nil, errInvalidLoggerInstance
	}
}

// NewLoggerWithModuleName returns an instance of logger
func NewLoggerWithModuleName(config *Configuration, moduleName string) (Logger, error) {
	logger, err := NewLogger(config, InstanceZapLogger)
	if err != nil {
		return nil, err
	}

	logger = logger.WithFields(Fields{
		"module": moduleName,
	})

	return logger, nil
}
