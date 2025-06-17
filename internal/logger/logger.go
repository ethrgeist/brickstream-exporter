package logger

import (
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"time"
)

var log zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	log = zerolog.New(output).With().Timestamp().Caller().Logger()

	if logLevelStr, exists := os.LookupEnv("GO_LOG"); exists {
		if logLevel, err := strconv.Atoi(logLevelStr); err == nil {
			switch logLevel {
			case -1:
				zerolog.SetGlobalLevel(zerolog.Disabled)
			case 0:
				zerolog.SetGlobalLevel(zerolog.PanicLevel)
			case 1:
				zerolog.SetGlobalLevel(zerolog.FatalLevel)
			case 2:
				zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			case 3:
				zerolog.SetGlobalLevel(zerolog.WarnLevel)
			case 4:
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			case 5:
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			default:
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
		}
	}
}

func GetLogger() zerolog.Logger {
	return log
}
