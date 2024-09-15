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

func NewLogger(level int, debugOn bool, logFilePath string) *Logger {
	l := &Logger{Level: level, DebugOn: debugOn}
	dir := filepath.Dir(logFilePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic("Failed to create log directory: " + err.Error())
	}
	l.LogFile, err = os.Create(logFilePath + "log.log")
	if err != nil {
		panic("Failed to create log file: " + err.Error())
	}
	l.Debug("Created log file", 0)
	return l

}

func (l *Logger) Debug(log string, level int) {
	if level >= l.Level && l.DebugOn {
		l.LogFile.WriteString(log + "\n")
	}
}
