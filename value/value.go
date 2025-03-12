package value

import (
	"context"
)

type Valuer func(ctx context.Context) any

// Value return the function value.
func Value(ctx context.Context, v any) any {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

func BindValues(ctx context.Context, keyvals []any) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v(ctx)
		}
	}
}

func ContainsValuer(keyvals []any) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}
