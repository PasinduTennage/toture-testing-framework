package util

import "os"

type Logger struct {
	Level   int
	DebugOn bool
	LogFile *os.File
}

func NewLogger(level int, debugOn bool, logFile string) *Logger {
	l := &Logger{Level: level, DebugOn: debugOn}
	if logFile != "" {
		l.LogFile, _ = os.Create(logFile)
		l.LogFile.WriteString("Log file created\n")
	} else {
		panic("Log file not specified")
	}
	return l

}

func (l *Logger) Debug(log string, level int) {
	if level >= l.Level && l.DebugOn {
		l.LogFile.WriteString(log)
	}
}
