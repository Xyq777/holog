package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ncuhome/holog"
	"github.com/ncuhome/holog/middleware/hogin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func TestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		span := trace.SpanFromContext(c.Request.Context())
		traceID := span.SpanContext().TraceID().String()

		c.String(200, "TraceID in test middleware: %s", traceID)

		c.Next()
	}
}

func initTracer() {

	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint("localhost:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		panic(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("your-service"),
		)),
	)
	otel.SetTracerProvider(tp)
}

func main() {
	initTracer()
	// logger := holog.NewLogger("test-service", holog.WithFileWriter(&lumberjack.Logger{
	// 	Filename:   "./zap.log",
	// 	MaxSize:    10,
	// 	MaxBackups: 5,
	// 	MaxAge:     30,
	// 	Compress:   false,
	// }))

	// logger.Error("errrrrrrrrrrrr")

	// logger.Info("www")
	r := gin.New()

	r.Use(hogin.Trace(), hogin.Logger(), TestMiddleware())

	r.GET("/", func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		traceID := span.SpanContext().TraceID().String()
		logger := holog.FromGinContext(c)
		logger.Info("12345")

		c.String(200, "TraceID: %s", traceID)
	})

	r.Run(":8080")
}
