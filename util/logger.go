package util

type Logger struct {
	Level   int
	DebugOn bool
}

func NewLogger(level int, debugOn bool) *Logger {
	return &Logger{Level: level, DebugOn: debugOn}
}

func (l *Logger) Debug(log string, level int) {
	if level >= l.Level && l.DebugOn {
		println(log)
	}
}
