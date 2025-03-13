package holog

func HError(err error) []any {
	return []any{"error", err}
}
