package tracing

import (
	"context"

	"github.com/ncuhome/holog/value"
	"go.opentelemetry.io/otel/trace"
)

func TraceID() value.Valuer {
	return func(ctx context.Context) any {
		if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.HasTraceID() {
			return spanCtx.TraceID().String()
		}
		return ""
	}
}

func SpanID() value.Valuer {
	return func(ctx context.Context) any {
		if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
			return span.SpanID().String()
		}
		return ""
	}
}
