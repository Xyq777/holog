package holog

import "sync"

var (
	globalLogger     *logger
	globalLoggerOnce sync.Once
	globalMu         sync.RWMutex
)

func init() {
	globalLoggerOnce.Do(func() {
		globalLogger = NewLogger("")
	})
}
func getGlobal() *logger {
	if globalLogger == nil {
		panic("logger not initialized")
	}
	return globalLogger
}

func SetGlobal(newLogger *logger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	old := globalLogger
	globalLogger = newLogger
	old.Close()
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
