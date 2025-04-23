package holog

import (
	"context"
	"fmt"
	"sync"
)

var (
	globalLogger     *Logger
	globalLoggerOnce sync.Once
	globalMu         sync.RWMutex
)

func init() {
	globalLoggerOnce.Do(func() {
		globalLogger = NewLogger("")
	})
}
func getGlobal() *Logger {
	if globalLogger == nil {
		panic("Logger not initialized")
	}
	return globalLogger
}

func SetGlobal(newLogger *Logger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	old := globalLogger
	globalLogger = newLogger
	old.Close()
}

func GetGlobal() *Logger {
	return getGlobal()
}

// Global methos
func Info(msg string, kvs ...any) {
	getGlobal().Info(msg, kvs...)
}
func Debug(msg string, kvs ...any) {
	getGlobal().Debug(msg, kvs...)
}
func Warn(msg string, kvs ...any) {
	getGlobal().Warn(msg, kvs...)
}
func Error(msg string, kvs ...any) {
	getGlobal().Error(msg, kvs...)
}

func Fatal(msg string, kvs ...any) {
	getGlobal().Fatal(msg, kvs...)
}
func Panic(msg string, kvs ...any) {
	getGlobal().Panic(msg, kvs...)
}

func Infof(format string, args ...any) {
	getGlobal().Info(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...any) {
	getGlobal().Debug(fmt.Sprintf(format, args...))
}
func Warnf(format string, args ...any) {
	getGlobal().Warn(fmt.Sprintf(format, args...))
}
func Errorf(format string, args ...any) {
	getGlobal().Error(fmt.Sprintf(format, args...))
}
func Fatalf(format string, args ...any) {
	getGlobal().Fatal(fmt.Sprintf(format, args...))
}
func Panicf(format string, args ...any) {
	getGlobal().Panic(fmt.Sprintf(format, args...))
}

func Ctx(ctx context.Context) *Logger {
	logger := getGlobal().copy()
	logger.ctx = ctx
	return logger
}
