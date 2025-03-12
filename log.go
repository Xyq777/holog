package holog

import (
	"context"

	"github.com/natefinch/lumberjack"
	"github.com/ncuhome/holog/level"
	"github.com/ncuhome/holog/value"
	"github.com/ncuhome/holog/zapLogger"
)

type Mode uint8

const (
	Dev Mode = iota
	Prod
)

type Logger interface {
	Log(level level.Level, msg string, kvs ...any)
	Close()
}

type logger struct {
	logger    Logger
	prefix    []any
	hasValuer bool
	ctx       context.Context
}

func NewLogger(serviceName string, opts ...Option) *logger {
	options := options{
		lumberjackLogger: nil,
		mode:             Prod,
	}
	for _, opt := range opts {
		opt(&options)
	}
	prefix := []any{
		"service", serviceName,
	}
	return &logger{logger: zapLogger.NewZappLogger(options.lumberjackLogger, uint8(options.mode)),
		prefix:    prefix,
		hasValuer: value.ContainsValuer(prefix),
		ctx:       context.Background()}
}

type Option func(o *options)

// Only supports zap now
type options struct {
	// logger Logger
	mode             Mode
	lumberjackLogger *lumberjack.Logger
}

func WithFileWriter(lumberjackLogger *lumberjack.Logger) Option {
	return func(o *options) {
		o.lumberjackLogger = lumberjackLogger
	}
}

func WithMode(mode Mode) Option {
	return func(o *options) {
		o.mode = mode
	}
}

func (l *logger) Close() {
	l.logger.Close()
}

func (l *logger) Info(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	l.logger.Log(level.InfoLevel, msg, getKeyVals(l.prefix, kvs)...)
}
func (l *logger) Warn(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	l.logger.Log(level.WarnLevel, msg, getKeyVals(l.prefix, kvs)...)
}
func (l *logger) Debug(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	l.logger.Log(level.DebugLevel, msg, getKeyVals(l.prefix, kvs)...)
}
func (l *logger) Error(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	l.logger.Log(level.ErrorLevel, msg, getKeyVals(l.prefix, kvs)...)
}
func (l *logger) Fatal(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	l.logger.Log(level.FatalLevel, msg, getKeyVals(l.prefix, kvs)...)
}
func (l *logger) Panic(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	l.logger.Log(level.PanicLevel, msg, getKeyVals(l.prefix, kvs)...)
}

func getKeyVals(prefix []any, kvs []any) []any {
	keyvals := make([]any, 0, len(prefix)+len(kvs))
	keyvals = append(keyvals, prefix...)
	keyvals = append(keyvals, kvs...)
	return keyvals
}
