package level

type Level uint8

const (
	InfoLevel Level = iota
	DebugLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

func (l *Level) ToString() string {
	switch *l {
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	default:
		return "NO_SUCH_LEVEL"
	}
}
