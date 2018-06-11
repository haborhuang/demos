package main

import (
	"go.uber.org/zap"
	"time"
	"gopkg.in/natefinch/lumberjack.v2"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	lj := lumberjack.Logger{
		Filename: "/Users/huangxin/go/src/foo/tmp/foo.log",
		MaxSize: 1,
		MaxBackups: 2,
		MaxAge: 28,
		LocalTime: true,
		Compress: true,
	}
	w := zapcore.AddSync(&lj)

	core :=  zapcore.NewTee(
		// rotated logger
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), w, zap.InfoLevel),
		// stdout logger
		zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.Lock(os.Stdout), zap.InfoLevel),
	)
	logger := zap.New(core)
	defer logger.Sync()


	go func() {
		for {
			time.Sleep(5 * time.Second)
			lj.Rotate()
		}
	}()

	for {
		time.Sleep(time.Second)

		logger.Info(
			"hello world",
			zap.Time("at", time.Now()),
		)
	}
}
