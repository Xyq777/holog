package holog

import (
	"context"
	"fmt"

	"github.com/natefinch/lumberjack"

	"github.com/ncuhome/holog/level"
	"github.com/ncuhome/holog/sink"
	"github.com/ncuhome/holog/tracing"
	"github.com/ncuhome/holog/value"
	"github.com/ncuhome/holog/zapLogger"
)

type Mode uint8

const (
	Dev Mode = iota
	Prod
)

type OutputStyle uint8

const (
	JSON OutputStyle = iota
	TEXT
)

type Logger interface {
	Log(level level.Level, msg string, kvs ...any) (sink.LogEntry, error)
	Close()
}

type logger struct {
	logger    Logger
	prefix    []any
	hasValuer bool
	ctx       context.Context
	mode      Mode
	sink      sink.Sink
}

func NewLogger(serviceName string, opts ...Option) *logger {
	options := options{
		lumberjackLogger: nil,
		mode:             Dev,
		style:            JSON,
		fields:           []any{},
		sink:             nil,
	}
	for _, opt := range opts {
		opt(&options)
	}
	prefix := []any{
		"service", serviceName,
		"timestamp", value.DefaultTimestamp,
		"caller", value.DefaultCaller,
		"trace_id", tracing.TraceID(),
	}
	prefix = append(prefix, options.fields...)
	return &logger{logger: zapLogger.NewZappLogger(options.lumberjackLogger, uint8(options.style)),
		prefix:    prefix,
		hasValuer: value.ContainsValuer(prefix),
		ctx:       context.Background(),
		mode:      options.mode,
		sink:      options.sink}
}

type Option func(o *options)

// Only supports zap now
type options struct {
	// logger Logger
	mode             Mode
	style            OutputStyle
	lumberjackLogger *lumberjack.Logger
	fields           []any
	sink             sink.Sink
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

func WithOutputStyle(style OutputStyle) Option {
	return func(o *options) {
		o.style = style
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

func WithSink(sink sink.Sink) Option {
	return func(o *options) {
		o.sink = sink
	}
}

func (l *logger) Close() {
	l.logger.Close()
}

func (l *logger) copy() *logger {
	return &logger{
		logger:    l.logger,
		prefix:    l.prefix,
		hasValuer: l.hasValuer,
		ctx:       l.ctx,
		mode:      l.mode,
		sink:      l.sink,
	}
}

func (l *logger) Info(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.InfoLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *logger) Warn(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.WarnLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *logger) Debug(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.DebugLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *logger) Error(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.ErrorLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *logger) Fatal(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.FatalLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *logger) Panic(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.PanicLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}

func getKeyVals(prefix []any, kvs []any) []any {
	keyvals := make([]any, 0, len(prefix)+len(kvs))
	keyvals = append(keyvals, prefix...)
	keyvals = append(keyvals, kvs...)
	return keyvals
}
