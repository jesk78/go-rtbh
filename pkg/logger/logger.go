package logger

func New(timestamp, debug bool) *Logger {
	return &Logger{
		UseTimestamp: timestamp,
		UseDebug:     debug,
	}
}
