package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logconf := zap.NewProductionConfig()
	logconf.DisableStacktrace = true
	logconf.Encoding = "console"
	logconf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, err := logconf.Build()
	if err != nil {
		panic(err)
	}

	app := NewApp(logger, &logconf.Level)
	if err := app.Command().Run(os.Args); err != nil {
		app.logger.Error("runtime error", zap.Error(err))
	}
}
