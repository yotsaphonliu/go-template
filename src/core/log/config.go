package log

import (
	"github.com/spf13/viper"
)

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
	Color             bool
}

func InitConfig() (*Configuration, error) {
	logLevel := viper.GetString("Log.Level")
	logColor := viper.GetBool("Log.Color")
	logJSON := viper.GetBool("Log.JSON")

	logLevel = NormalizeLogLevel(logLevel)

	config := &Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logLevel,
		Color:             logColor,
		ConsoleJSONFormat: logJSON,
	}

	return config, nil
}

func NormalizeLogLevel(logLevel string) string {
	var normalizedLogLevel string
	switch logLevel {
	case "info":
		normalizedLogLevel = Info
	case "debug":
		normalizedLogLevel = Debug
	case "warn":
		normalizedLogLevel = Warn
	case "error":
		normalizedLogLevel = Error
	case "fatal":
		normalizedLogLevel = Fatal
	default:
		normalizedLogLevel = Debug
	}
	return normalizedLogLevel
}
