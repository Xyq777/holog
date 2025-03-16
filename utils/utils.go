package utils

import (
	"fmt"

	"github.com/ncuhome/holog/sink"
)

func DataToLogEntry(data []any) (sink.LogEntry, error) {
	result := make(map[string]interface{})
	for i := 0; i < len(data); i += 2 {
		key, ok := data[i].(string)
		if !ok {
			return nil, fmt.Errorf("key is not a string")
		}
		value := data[i+1]
		result[key] = value
	}
	return result, nil
}
