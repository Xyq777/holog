package ingester

import (
	"context"
)

type LogEntry map[string]interface{}

type Ingester interface {
	Send(ctx context.Context, entry LogEntry) error
	SendBatch(ctx context.Context, entries []LogEntry) error
}
