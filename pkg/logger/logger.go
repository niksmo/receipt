package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

const defaultLevel = "info"

type Logger struct {
	zerolog.Logger
}

func New(level string) Logger {
	setLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	return Logger{
		zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

func setLevel(level string) {
	const op = "logger.setLevel"

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
		fmt.Printf(
			"unexpected log level: %q\nset to default: %q\n",
			level, defaultLevel,
		)
	}
	zerolog.SetGlobalLevel(lvl)
}

func (l Logger) WithOp(op string) Logger {
	return Logger{l.With().Str("op", op).Logger()}
}
