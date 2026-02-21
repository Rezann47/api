package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New uygulama ortamına göre uygun logger döner
func New(env string) *zap.Logger {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	log, err := cfg.Build(zap.AddCallerSkip(0))
	if err != nil {
		panic(err)
	}
	return log
}
