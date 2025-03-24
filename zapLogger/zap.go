package zapLogger

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/ncuhome/holog/level"
	"github.com/ncuhome/holog/sink"
	"github.com/ncuhome/holog/utils"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log *zap.Logger
}

var DefaultEncoder = zapcore.EncoderConfig{
	LevelKey:       "level",
	NameKey:        "logger",
	MessageKey:     "message",
	StacktraceKey:  "stack",
	EncodeCaller:   zapcore.FullCallerEncoder,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
}

func NewZappLogger(lumberjackLogger *lumberjack.Logger, exporter *otlploghttp.Exporter, serviceName string, mode uint8) *ZapLogger {
	return newZapLoggerWithConfigs(
		DefaultEncoder,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		lumberjackLogger,
		exporter,
		mode,
		serviceName,
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

func newZapLoggerWithConfigs(encoder zapcore.EncoderConfig, level zap.AtomicLevel, lumberJackLogger *lumberjack.Logger, exporter *otlploghttp.Exporter, style uint8, serviceName string, opts ...zap.Option) *ZapLogger {
	level.SetLevel(zap.InfoLevel)
	var core zapcore.Core
	var processor *log.BatchProcessor
	var provider *log.LoggerProvider
	if exporter != nil {
		processor = log.NewBatchProcessor(exporter)
		provider = log.NewLoggerProvider(
			log.WithProcessor(processor),
		)
	} else {
		provider = log.NewLoggerProvider()
	}
	if lumberJackLogger != nil {
		writeSyncer := getLogWriter(lumberJackLogger)
		if style != 1 {
			core = zapcore.NewTee(zapcore.NewCore(
				zapcore.NewJSONEncoder(encoder),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)),
				level,
			), otelzap.NewCore(serviceName, otelzap.WithLoggerProvider(provider)))
		} else {
			core = zapcore.NewTee(zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoder),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)),
				level,
			), otelzap.NewCore(serviceName, otelzap.WithLoggerProvider(provider)))
		}

	} else {
		if style != 1 {
			core = zapcore.NewTee(zapcore.NewCore(
				zapcore.NewJSONEncoder(encoder),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
				level,
			), otelzap.NewCore(serviceName, otelzap.WithLoggerProvider(provider)))
		} else {
			core = zapcore.NewTee(zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoder),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
				level,
			), otelzap.NewCore(serviceName, otelzap.WithLoggerProvider(provider)))
		}
	}
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger}
}
func (logger *ZapLogger) Log(l level.Level, msg string, kvs ...any) (sink.LogEntry, error) {

	if len(kvs)%2 != 0 {
		logger.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", kvs))
		return nil, fmt.Errorf("keyvalues must appear in pairs: %v", kvs)
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
	return utils.DataToLogEntry(kvs)
}

func (logger *ZapLogger) Close() {
	// TODO?
	if err := logger.log.Sync(); err != nil {

	}
}
