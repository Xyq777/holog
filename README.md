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
```golang
// 配置文件输出选项，可自定义文件输出配置，若无该选项则默认不输出到文件
WithFileWriter(lumberjackLogger *lumberjack.Logger)
// 配置日志模式，Dev模式会以TEXT格式输出，Prod模式会以JSON格式输出，如无该选项则默认以JSON格式输出
WithMode(mode Mode)
```
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

