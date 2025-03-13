package holog

import (
	"context"
	"fmt"

	"github.com/natefinch/lumberjack"

	"github.com/ncuhome/holog/ingester"
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
	Log(level level.Level, msg string, kvs ...any) (ingester.LogEntry, error)
	Close()
}

type logger struct {
	logger    Logger
	prefix    []any
	hasValuer bool
	ctx       context.Context
	mode      Mode
	ingester  ingester.Ingester
}

func NewLogger(serviceName string, opts ...Option) *logger {
	options := options{
		lumberjackLogger: nil,
		mode:             Prod,
		fields:           []any{},
		ingester:         ingester.NewO2Imgester(),
	}
	for _, opt := range opts {
		opt(&options)
	}
	prefix := []any{
		"service", serviceName,
		"timestamp", value.DefaultTimestamp,
	}
	prefix = append(prefix, options.fields...)
	return &logger{logger: zapLogger.NewZappLogger(options.lumberjackLogger, uint8(options.mode)),
		prefix:    prefix,
		hasValuer: value.ContainsValuer(prefix),
		ctx:       context.Background(),
		mode:      options.mode,
		ingester:  options.ingester}
}

type Option func(o *options)

// Only supports zap now
type options struct {
	// logger Logger
	mode             Mode
	lumberjackLogger *lumberjack.Logger
	fields           []any
	ingester         ingester.Ingester
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

func WithFields(fields ...any) Option {
	if len(fields) != 0 && len(fields)%2 != 0 {
		panic(fmt.Sprintf("Keyvalues must appear in pairs: %v", fields...))
	}
	return func(o *options) {
		o.fields = fields
	}
}

func WithIngester(ingester ingester.Ingester) Option {
	return func(o *options) {
		o.ingester = ingester
	}
}

func (l *logger) Close() {
	l.logger.Close()
}

func (l *logger) Info(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	logEntry, err := l.logger.Log(level.InfoLevel, msg, getKeyVals(l.prefix, kvs)...)
	if l.ingester != nil && err != nil && l.mode == Prod {
		l.ingester.Send(l.ctx, logEntry)
	}
}
func (l *logger) Warn(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	logEntry, err := l.logger.Log(level.WarnLevel, msg, getKeyVals(l.prefix, kvs)...)
	if l.ingester != nil && err != nil && l.mode == Prod {
		l.ingester.Send(l.ctx, logEntry)
	}
}
func (l *logger) Debug(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	logEntry, err := l.logger.Log(level.DebugLevel, msg, getKeyVals(l.prefix, kvs)...)
	if l.ingester != nil && err != nil && l.mode == Prod {
		l.ingester.Send(l.ctx, logEntry)
	}
}
func (l *logger) Error(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	logEntry, err := l.logger.Log(level.ErrorLevel, msg, getKeyVals(l.prefix, kvs)...)
	if l.ingester != nil && err != nil && l.mode == Prod {
		l.ingester.Send(l.ctx, logEntry)
	}
}
func (l *logger) Fatal(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	logEntry, err := l.logger.Log(level.FatalLevel, msg, getKeyVals(l.prefix, kvs)...)
	if l.ingester != nil && err != nil && l.mode == Prod {
		l.ingester.Send(l.ctx, logEntry)
	}
}
func (l *logger) Panic(msg string, kvs ...any) {
	if l.hasValuer {
		value.BindValues(l.ctx, l.prefix)
	}
	logEntry, err := l.logger.Log(level.PanicLevel, msg, getKeyVals(l.prefix, kvs)...)
	if l.ingester != nil && err != nil && l.mode == Prod {
		l.ingester.Send(l.ctx, logEntry)
	}
}

func getKeyVals(prefix []any, kvs []any) []any {
	keyvals := make([]any, 0, len(prefix)+len(kvs))
	keyvals = append(keyvals, prefix...)
	keyvals = append(keyvals, kvs...)
	return keyvals
}
