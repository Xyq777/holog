package hogin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ncuhome/holog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	epoch     = time.Date(2025, 3, 18, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6
	nodeID    = int64(1)
	mutex     sync.Mutex
	lastStamp int64
	counter   int64
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := otel.GetTracerProvider()
		_, isSDKProvider := provider.(*sdktrace.TracerProvider)

		ctx := propagation.TraceContext{}.Extract(
			c.Request.Context(),
			propagation.HeaderCarrier(c.Request.Header),
		)

		if !isSDKProvider {
			ctx = withSnowflakeTraceID(ctx)
		} else {
			tracer := provider.Tracer("gin")
			var span trace.Span
			ctx, span = tracer.Start(ctx, c.Request.URL.Path)
			defer span.End()
		}

		c.Request = c.Request.WithContext(ctx)
		propagation.TraceContext{}.Inject(ctx, propagation.HeaderCarrier(c.Writer.Header()))

		logger := holog.CopyGlobalWithContext(ctx)
		c.Set("logger", logger)
		c.Next()
	}
}

func snowflakeTraceID() string {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now().UnixNano()/1e6 - epoch
	if now == lastStamp {
		counter++
	} else {
		counter = 0
		lastStamp = now
	}

	id := (now << 22) | (nodeID << 12) | (counter % 4096)
	b := make([]byte, 16)
	rand.Read(b[8:])

	b[0] = byte(id >> 56)
	b[1] = byte(id >> 48)
	b[2] = byte(id >> 40)
	b[3] = byte(id >> 32)
	b[4] = byte(id >> 24)
	b[5] = byte(id >> 16)
	b[6] = byte(id >> 8)
	b[7] = byte(id)

	return hex.EncodeToString(b)
}

func withSnowflakeTraceID(ctx context.Context) context.Context {
	traceID := snowflakeTraceID()
	tid, _ := trace.TraceIDFromHex(traceID)
	sid := generateSpanID()

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(ctx, spanCtx)
}

func generateSpanID() trace.SpanID {
	var sid [8]byte
	rand.Read(sid[:])
	return sid
}
