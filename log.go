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
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
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

type LoggerItf interface {
	Log(level level.Level, msg string, kvs ...any) (sink.LogEntry, error)
	Close()
}

type Logger struct {
	logger    LoggerItf
	prefix    []any
	hasValuer bool
	ctx       context.Context
	mode      Mode
	sink      sink.Sink
}

func NewLogger(serviceName string, opts ...Option) *Logger {
	options := options{
		lumberjackLogger: nil,
		mode:             Dev,
		style:            JSON,
		fields:           []any{},
		sink:             nil,
		exporter:         nil,
	}
	for _, opt := range opts {
		opt(&options)
	}
	prefix := []any{
		"service", serviceName,
		"timestamp", value.DefaultTimestamp,
		"caller", value.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	}
	prefix = append(prefix, options.fields...)
	return &Logger{logger: zapLogger.NewZappLogger(options.lumberjackLogger, options.exporter, serviceName, uint8(options.style)),
		prefix:    prefix,
		hasValuer: value.ContainsValuer(prefix),
		ctx:       context.Background(),
		mode:      options.mode,
		sink:      options.sink}
}

type Option func(o *options)

// Only supports zap now
type options struct {
	// Logger LoggerItf
	mode             Mode
	style            OutputStyle
	lumberjackLogger *lumberjack.Logger
	fields           []any
	sink             sink.Sink
	exporter         *otlploghttp.Exporter
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

func WithExporter(exporter *otlploghttp.Exporter) Option {
	return func(o *options) {
		o.exporter = exporter
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

func (l *Logger) Close() {
	l.logger.Close()
}

func (l *Logger) copy() *Logger {
	return &Logger{
		logger:    l.logger,
		prefix:    l.prefix,
		hasValuer: l.hasValuer,
		ctx:       l.ctx,
		mode:      l.mode,
		sink:      l.sink,
	}
}

func (l *Logger) Info(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.InfoLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *Logger) Warn(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.WarnLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *Logger) Debug(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.DebugLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *Logger) Error(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.ErrorLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *Logger) Fatal(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.FatalLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}
func (l *Logger) Panic(msg string, kvs ...any) {
	keyvals := getKeyVals(l.prefix, kvs)
	if l.hasValuer {
		value.BindValues(l.ctx, keyvals)
	}
	logEntry, err := l.logger.Log(level.PanicLevel, msg, keyvals...)
	if l.sink != nil && err != nil && l.mode == Prod {
		l.sink.Send(l.ctx, logEntry)
	}
}

func (l *Logger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.Warn(fmt.Sprintf(format, args...))
}
func (l *Logger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}
func (l *Logger) Fatalf(format string, args ...any) {
	l.Fatal(fmt.Sprintf(format, args...))
}

func (l *Logger) Panicf(format string, args ...any) {
	l.Panic(fmt.Sprintf(format, args...))
}

func (l *Logger) Ctx(ctx context.Context) *Logger {
	logger := l.copy()
	logger.ctx = ctx
	return logger
}

func getKeyVals(prefix []any, kvs []any) []any {
	keyvals := make([]any, 0, len(prefix)+len(kvs))
	keyvals = append(keyvals, prefix...)
	keyvals = append(keyvals, kvs...)
	return keyvals
}
