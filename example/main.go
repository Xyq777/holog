package main

import (
	"errors"

	"github.com/natefinch/lumberjack"
	"github.com/ncuhome/holog"
)

func main() {
	logger := holog.NewLogger("test-service", holog.WithFileWriter(&lumberjack.Logger{
		Filename:   "./zap.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}))

	// logger := holog.NewLogger("test-service", holog, holog.WithMode(holog.Dev), holog.WithFileWriter(&lumberjack.Logger{
	// 	Filename:   "./zap.log",
	// 	MaxSize:    10,
	// 	MaxBackups: 5,
	// 	MaxAge:     30,
	// 	Compress:   false,
	// }))

	logger.Info("This is a test info with message")
	logger.Info("")
	logger.Info("This is a test info with message and kvs", "code", 200)

	holog.Error("This is a test error in default global holog logger")

	holog.SetGlobal(logger)
	holog.Error("This is a test error in customized global holog logger")

	err := errors.New("test error")
	holog.Error("This is a test error", "error", err)

}
