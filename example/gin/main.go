package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ncuhome/holog"
	"github.com/ncuhome/holog/middleware/hogin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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
			semconv.ServiceNameKey.String("test-service"),
		)),
	)
	otel.SetTracerProvider(tp)
}

func main() {
	initTracer()
	r := gin.New()
	holog.Info("haha")
	r.Use(otelgin.Middleware("test-service"), hogin.Logger())
	r.GET("/", func(c *gin.Context) {
		holog.Ctx(c.Request.Context()).Info("hahaha")
	})
	r.Run(":8080")
}
