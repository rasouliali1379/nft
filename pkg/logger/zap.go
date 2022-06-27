package logger

import (
	"context"
	"log"
	"maskan/client/elk"
	"maskan/config"
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitGlobalLogger(lc fx.Lifecycle, logstash *elk.LogStash) error {
	logger := configLogger(logstash)
	zap.ReplaceGlobals(logger)

	lc.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			if err := zap.L().Sync(); err != nil {
				log.Println("logger failed to sync:", err)
			}
			return nil
		},
	})
	return nil
}

func configLogger(el *elk.LogStash) *zap.Logger {
	logLevel := getLogLevel()

	elkZapCore := ecszap.NewCore(
		ecszap.NewDefaultEncoderConfig(),
		zapcore.AddSync(el),
		logLevel,
	)

	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	terminalZapCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout),
		logLevel,
	)

	core := zapcore.NewTee(elkZapCore, terminalZapCore)
	logger := zap.New(core, zap.AddCaller())
	return logger.With(zap.String("service", config.C().App.Name))
}

func getLogLevel() zapcore.Level {
	if config.C().App.DebugMode {
		return zap.DebugLevel
	}
	return zap.InfoLevel
}
