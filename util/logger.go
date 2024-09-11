package util

import (
	"os"
	"path/filepath"
)

type Logger struct {
	Level   int
	DebugOn bool
	LogFile *os.File
}

func NewLogger(level int, debugOn bool, logFile string) *Logger {
	l := &Logger{Level: level, DebugOn: debugOn}
	if logFile != "" {

		dir := filepath.Dir(logFile)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic("Failed to create log directory: " + err.Error())
		}
		l.LogFile, err = os.Create(logFile)
		if err != nil {
			panic("Failed to create log file: " + err.Error())
		}
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
