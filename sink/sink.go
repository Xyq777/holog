package sink

import "context"

type Sink interface {
	Send(ctx context.Context, entry LogEntry) error
	SendBatch(ctx context.Context, entries []LogEntry) error
}

type LogEntry map[string]interface{}
