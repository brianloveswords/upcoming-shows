package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Logger struct {
	Stack       []string
	Stderr      io.Writer
	Stdout      io.Writer
	DebugPrefix string
	Level       LoggerLevel
}

type LoggerLevel int

const (
	LoggerLevelSilent LoggerLevel = iota
	LoggerLevelNormal
	LoggerLevelVerbose
	LoggerLevelDebug
	LoggerLevelExtreme
)

var DefaultLoggerDebugPrefix = "[debug] "

func NewLogger() Logger {
	return Logger{
		Stderr:      os.Stderr,
		Stdout:      os.Stdout,
		DebugPrefix: DefaultLoggerDebugPrefix,
		Level:       LoggerLevelNormal,
	}
}
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Log(format, v...)
	os.Exit(1)
}

func (l *Logger) Extreme(format string, v ...interface{}) {
	if l.Level < LoggerLevelExtreme {
		return
	}
	l.Debug(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.Level < LoggerLevelDebug {
		return
	}
	fpath := l.stackString()
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(l.Stderr, l.DebugPrefix+fpath+": "+msg+"\n")
}
func (l *Logger) Verbose(format string, v ...interface{}) {
	if l.Level < LoggerLevelVerbose {
		return
	}
	l.stderrPrint(format, v...)
}
func (l *Logger) Normal(format string, v ...interface{}) {
	l.Log(format, v...)
}
func (l *Logger) Log(format string, v ...interface{}) {
	if l.Level < LoggerLevelNormal {
		return
	}
	l.stderrPrint(format, v...)
}
func (l *Logger) Prompt(format string, v ...interface{}) {
	if l.Level < LoggerLevelNormal {
		return
	}
	fmt.Fprintf(l.Stderr, format+": ", v...)
}
func (l *Logger) stderrPrint(format string, v ...interface{}) {
	fmt.Fprintf(l.Stderr, format+"\n", v...)
}

func (l *Logger) CmdOutput(format string, v ...interface{}) {
	l.stdoutPrint(format, v...)
}
func (l *Logger) stdoutPrint(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Fprint(l.Stdout, msg+"\n")
}
func (l *Logger) Enter(name string) func() {
	l.Stack = append(l.Stack, name)
	l.Extreme("entering %s", name)
	return func() {
		if l.Level >= LoggerLevelExtreme {
			l.Extreme("exiting %s", name)
		}
		l.Stack = l.Stack[:len(l.Stack)-1]
	}
}

func (l *Logger) IsLevelSilent() bool  { return l.Level <= LoggerLevelSilent }
func (l *Logger) IsLevelNormal() bool  { return l.Level <= LoggerLevelNormal }
func (l *Logger) IsLevelLog() bool     { return l.Level <= LoggerLevelNormal }
func (l *Logger) IsLevelVerbose() bool { return l.Level <= LoggerLevelVerbose }
func (l *Logger) IsLevelDebug() bool   { return l.Level <= LoggerLevelDebug }
func (l *Logger) IsLevelExtreme() bool { return l.Level <= LoggerLevelExtreme }

func (l *Logger) stackString() string {
	fpath := strings.Join(l.Stack, ":")
	if fpath == "" {
		fpath = "<top>"
	}
	return fpath
}
