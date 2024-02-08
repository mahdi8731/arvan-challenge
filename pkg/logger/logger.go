package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(levelStr string) *zerolog.Logger {
	level, err := zerolog.ParseLevel(levelStr)

	if err != nil {
		panic("invalid logger level!!!")
	}
	zerolog.SetGlobalLevel(level)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	return &logger
}
