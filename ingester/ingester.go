package ingester

import (
	"context"
)

type LogEntry map[string]interface{}

type Ingester interface {
	Send(ctx context.Context, stream string, entry LogEntry) error
	SendBatch(ctx context.Context, stream string, entries []LogEntry) error
}
