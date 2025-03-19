package main

import (
	"context"
	"fmt"

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
	r := gin.New()
	holog.Info("少时诵诗书")
	holog.Info("少时诵诗书2")
	holog.Infof("%s菲", "ta")
	r.Use(hogin.Trace(), hogin.Logger())
	r.GET("/", func(c *gin.Context) {
		spanCtx := trace.SpanContextFromContext(c.Request.Context())
		fmt.Printf("spanCtx.TraceID(): %v\n", spanCtx.TraceID())
		logger := holog.FromGinContext(c)
		logger.Info("12345")

	})
	r.Run(":8080")
}
