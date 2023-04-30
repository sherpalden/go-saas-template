package infrastructure

import (
	"go.uber.org/zap"
)

type Logger struct {
	Zap *zap.SugaredLogger
}

func NewLogger() Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return Logger{
		Zap: sugar,
	}
}
