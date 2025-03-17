## Quick Start
#### $env:GOPRIVATE="github.com/ncuhome"
```shell
go get -u github.com/ncuhome/holog 
```
### 直接使用全局 *logger*
```golang
package main

import (
	"errors"

	"github.com/ncuhome/holog"
)


func main() {
    holog.Info("This is an info log")
    holog.Error("This is an error log")

    holog.Info("This is an info log with kvs","test key","test val")

    err:=errors.New("This is an error")
    holog.Error("This is an error log with error","error",err)
}

```
### 自定义 *logger*
当前可自定义选项有：

* WithFileWriter(lumberjackLogger *lumberjack.Logger)\
 说明：自定义输出文件配置。若不启用此选择，则不会将日志输出到文件里
* WithOutputStyle(style OutputStyle)\
 说明：自定义输出样式，有两种样式：*holog.JSON* 和 *holog.TEXT*。若不启动此选项，则默认为 *holog.JSON* 样式
* WithMode(mode Mode)\
 说明：自定义日志模式，有两种模式：*holog.Dev* 和 *holog.Prod*。*holog.Prod* 模式会将日志同步到日志系统（即 Sink 接口，默认为 *OpenObserve* ）。默认是 *holog.Dev*
* WithFields(fields ...any)\
 说明：添加自定义输出字段，至少传入偶数个参数，参数为 *key-value* 对，*key* 为 *string* 类型，*value* 为任意类型（后文会说明如何添加运行时变化的字段，如 *trace_id*）
* WithSink(sink Sink)\
 说明：自定义日志输出端，默认是 *nil* ， 只有 *Mode* 为 *holog.Prod* 时才会生效，则不会将日志输出到外部端，

```golang
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

	logger.Info("This is an info log")

	err := errors.New("This is an error")
	logger.Error("This is an error log with error", "error", err)
}
```
### 将自定义 *logger* 配置到全局
```golang
logger := holog.NewLogger("test-service", holog.WithFileWriter(&lumberjack.Logger{
	Filename:   "./zap.log",
	MaxSize:    10,
	MaxBackups: 5,
	MaxAge:     30,
	Compress:   false,
}))
holog.SetGlobal(logger)

holog.Info("This is a log from a new global logger")
```
### 自定义输出字段
加入普通字段：
```golang
logger := holog.NewLogger("test-service", holog.WithFields("new_key", "new_value"))
```
加入运行时变化的字段（如时间戳、trace_id）：
```golang
// 当前默认Valuer只有一个作为示例的DefaultTimestamp
logger := holog.NewLogger("test", holog.WithFields("ts", value.DefaultTimestamp))
```
运行时值 Valuer 是 func(context.Context)any 若要自定义运行时字段，请参照 value/value.go：
```golang
var (
	DefaultTimestamp = Timestamp(time.RFC3339)
)
func Timestamp(layout string) Valuer {
	return func(context.Context) any {
		return time.Now().Format(layout)
	}
}
```
## 接口
```golang
// 创建logger并使用
func NewLogger(serviceName string, opts ...Option) *logger
func (l *logger) Close()
func (l *logger) Info(msg string, kvs ...any) 
func (l *logger) Warn(msg string, kvs ...any) 
func (l *logger) Debug(msg string, kvs ...any)
func (l *logger) Error(msg string, kvs ...any) 
func (l *logger) Fatal(msg string, kvs ...any)
func (l *logger) Panic(msg string, kvs ...any)

// 自定义选项
func WithFileWriter(lumberjackLogger *lumberjack.Logger) Option 
func WithMode(mode Mode) Option 
func WithOutputStyle(style OutputStyle) Option 
func WithFields(fields ...any) Option 
func WithSink(sink sink.Sink) Option


// 创建logger后绑定到全局，方便使用
func SetGlobal(newLogger *logger)

func Info(msg string, kvs ...any) 
func Debug(msg string, kvs ...any)
func Warn(msg string, kvs ...any) 
func Error(msg string, kvs ...any)
func Fatal(msg string, kvs ...any)
func Panic(msg string, kvs ...any)

// 日志的外部输出端，可以是 OpenObserve、ElasticSearch、Kafka 等
type Sink interface {
	Send(ctx context.Context, entry LogEntry) error
	SendBatch(ctx context.Context, entries []LogEntry) error
}

```

