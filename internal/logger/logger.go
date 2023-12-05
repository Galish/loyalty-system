package logger

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

var logger = &log.Logger{
	Formatter: new(log.JSONFormatter),
	Out:       io.Discard,
}

type Fields map[string]interface{}

func Initialize(level string) {
	logger.SetOutput(os.Stderr)

	logLevel, err := log.ParseLevel(level)
	if err != nil {
		logger.Error("invalid log level", err)
		return
	}

	logger.SetLevel(logLevel)
}

func Debug(args ...interface{}) {
	logger.Log(log.DebugLevel, args...)
}

func Info(args ...interface{}) {
	logger.Log(log.InfoLevel, args...)
}

func WithError(err error) *log.Entry {
	return log.NewEntry(logger).WithError(err)
}

func WithFields(fields Fields) *log.Entry {
	f := log.Fields{}

	for k, v := range fields {
		f[k] = v
	}

	return log.NewEntry(logger).WithFields(f)
}
