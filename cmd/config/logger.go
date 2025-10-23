package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(isDebugMode bool) {
	var cfg zap.Config
	if isDebugMode {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	logger, _ := cfg.Build()

	logger.Info("Init Logger OK")

	zap.ReplaceGlobals(logger)
}
