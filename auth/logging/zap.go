package logging

import (
	"os"

	"go.uber.org/zap"
)

type zaplogger struct {
	Logger *zap.SugaredLogger
}

func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return &zaplogger{Logger: sugar}
}

func (z *zaplogger) Error(args ...interface{}) {
	z.Logger.Error(args...)
}

func (z *zaplogger) Warn(args ...interface{}) {
	z.Logger.Warn(args...)
}

func (z *zaplogger) Info(args ...interface{}) {
	z.Logger.Info(args...)
}

func (z *zaplogger) Fatal(args ...interface{}) {
	z.Logger.Fatal(args...)
	os.Exit(1)
}
