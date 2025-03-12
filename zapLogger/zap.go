package zapLogger

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/ncuhome/holog/level"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

var DefaultEncoder = zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stack",
	EncodeTime:     zapcore.RFC3339TimeEncoder,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.FullCallerEncoder,
}

func NewZappLogger(lumberjackLogger *lumberjack.Logger, mode uint8) *ZapLogger {
	return newZapLoggerWithConfigs(
		DefaultEncoder,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		lumberjackLogger,
		mode,
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
}

func getLogWriter(lumberJackLogger *lumberjack.Logger) zapcore.WriteSyncer {
	return zapcore.AddSync(lumberJackLogger)
}

func newZapLoggerWithConfigs(encoder zapcore.EncoderConfig, level zap.AtomicLevel, lumberJackLogger *lumberjack.Logger, mode uint8, opts ...zap.Option) *ZapLogger {
	level.SetLevel(zap.InfoLevel)
	var core zapcore.Core
	if lumberJackLogger != nil {
		writeSyncer := getLogWriter(lumberJackLogger)
		if mode != 0 {
			core = zapcore.NewCore(
				zapcore.NewJSONEncoder(encoder),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)),
				level,
			)
		} else {
			core = zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoder),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)),
				level,
			)
		}

	} else {
		if mode != 0 {
			core = zapcore.NewCore(
				zapcore.NewJSONEncoder(encoder),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
				level,
			)
		} else {
			core = zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoder),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
				level,
			)
		}
	}
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}
func (logger *ZapLogger) Log(l level.Level, msg string, kvs ...any) {

	if len(kvs)%2 != 0 {
		logger.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", kvs))
		return
	}
	var data []zap.Field

	for i := 0; i < len(kvs); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(kvs[i]), kvs[i+1]))
	}
	switch l {
	case level.InfoLevel:
		logger.log.Info(msg, data...)
	case level.DebugLevel:
		logger.log.Debug(msg, data...)
	case level.WarnLevel:
		logger.log.Warn(msg, data...)
	case level.ErrorLevel:
		logger.log.Error(msg, data...)
	case level.FatalLevel:
		logger.log.Fatal(msg, data...)
	case level.PanicLevel:
		logger.log.Panic(msg, data...)
	}
}

func (logger *ZapLogger) Close() {
	if err := logger.Sync(); err != nil {
		_, _ = os.Stderr.WriteString("failed to sync logger: " + err.Error() + "\n")
	}
}
