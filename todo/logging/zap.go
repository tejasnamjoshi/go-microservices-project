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

func (z *zaplogger) Error(err string) {
	z.Logger.Error(err)
}

func (z *zaplogger) Warn(warn ...string) {
	z.Logger.Warn(warn)
}
func (z *zaplogger) Info(info ...string) {
	z.Logger.Info(info)
}
func (z *zaplogger) Fatal(fatal string) {
	z.Logger.Fatal(fatal)
	os.Exit(1)
}
