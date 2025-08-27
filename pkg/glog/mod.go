package glog

import (
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger alias, just to ease a future possible database-migration to another glog library...
type Logger = zerolog.Logger

func InitLogger(level, service string, file *os.File) Logger {
	zerolog.TimestampFieldName = "date"
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	multi := zerolog.MultiLevelWriter(os.Stdout, file)

	logger := zerolog.New(multi).With().
		Timestamp().
		Logger()

	l, err := zerolog.ParseLevel(level)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error while setting glog level to %s", level)
	}

	logger = logger.Level(l)
	zerolog.SetGlobalLevel(l)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log.SetFlags(0)
	log.SetOutput(logger)

	logger = logger.With().Str(KeyService, service).Logger()

	zerolog.DefaultContextLogger = &logger

	return logger
}
