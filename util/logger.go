package util

type Logger struct {
	level   int
	debugOn bool
}

func NewLogger(level int, debugOn bool) *Logger {
	return &Logger{level: level, debugOn: debugOn}
}

func (l *Logger) Debug(log string, level int) {
	if l.level >= level && l.debugOn {
		println(log)
	}
}
