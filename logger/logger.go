package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Logger struct {
	Stack       []string
	Stderr      io.Writer
	Stdout      io.Writer
	DebugPrefix string
	Level       Level
}

type Level int

const (
	LevelSilent Level = iota
	LevelNormal
	LevelVerbose
	LevelDebug
	LevelExtreme
)

var DefaultDebugPrefix = "[debug] "

func New() Logger {
	return Logger{
		Stderr:      os.Stderr,
		Stdout:      os.Stdout,
		DebugPrefix: DefaultDebugPrefix,
		Level:       LevelNormal,
	}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	if l.Level == LevelSilent {
		return ioutil.Discard.Write(p)
	}
	return l.Stderr.Write(p)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Log(format, v...)
	os.Exit(1)
}

func (l *Logger) Extreme(format string, v ...interface{}) {
	if l.Level < LevelExtreme {
		return
	}
	l.Debug(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.Level < LevelDebug {
		return
	}
	fpath := l.stackString()
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(l.Stderr, l.DebugPrefix+fpath+": "+msg+"\n")
}
func (l *Logger) Verbose(format string, v ...interface{}) {
	if l.Level < LevelVerbose {
		return
	}
	l.stderrPrint(format, v...)
}
func (l *Logger) Normal(format string, v ...interface{}) {
	l.Log(format, v...)
}
func (l *Logger) Log(format string, v ...interface{}) {
	if l.Level < LevelNormal {
		return
	}
	l.stderrPrint(format, v...)
}
func (l *Logger) Prompt(format string, v ...interface{}) {
	if l.Level < LevelNormal {
		return
	}
	fmt.Fprintf(l.Stderr, format+": ", v...)
}
func (l *Logger) CmdOutput(format string, v ...interface{}) {
	l.stdoutPrint(format, v...)
}
func (l *Logger) Enter(name string) func() {
	l.Stack = append(l.Stack, name)
	l.Extreme("entering %s", name)
	return func() {
		if l.Level >= LevelExtreme {
			l.Extreme("exiting %s", name)
		}
		l.Stack = l.Stack[:len(l.Stack)-1]
	}
}

func (l *Logger) IsLevelSilent() bool  { return l.Level == LevelSilent }
func (l *Logger) IsLevelNormal() bool  { return l.Level >= LevelNormal }
func (l *Logger) IsLevelLog() bool     { return l.Level >= LevelNormal }
func (l *Logger) IsLevelVerbose() bool { return l.Level >= LevelVerbose }
func (l *Logger) IsLevelDebug() bool   { return l.Level >= LevelDebug }
func (l *Logger) IsLevelExtreme() bool { return l.Level >= LevelExtreme }

func (l *Logger) stdoutPrint(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Fprint(l.Stdout, msg+"\n")
}

func (l *Logger) stderrPrint(format string, v ...interface{}) {
	fmt.Fprintf(l.Stderr, format+"\n", v...)
}
func (l *Logger) stackString() string {
	fpath := strings.Join(l.Stack, ":")
	if fpath == "" {
		fpath = "<top>"
	}
	return fpath
}
func getLevelFromEnv() Level {
	switch os.Getenv("LOGLEVEL") {
	case "silent":
		return LevelSilent
	case "verbose":
		return LevelVerbose
	case "debug":
		return LevelDebug
	case "extreme":
		return LevelExtreme
	default:
		return LevelNormal
	}
}

var glog = Logger{
	Stderr:      os.Stderr,
	Stdout:      os.Stdout,
	DebugPrefix: DefaultDebugPrefix,
	Level:       getLevelFromEnv(),
}

var DefaultLogger = glog

var CmdOutput = glog.CmdOutput
var Debug = glog.Debug
var Enter = glog.Enter
var Extreme = glog.Extreme
var Fatal = glog.Fatal
var IsLevelDebug = glog.IsLevelDebug
var IsLevelExtreme = glog.IsLevelExtreme
var IsLevelLog = glog.IsLevelLog
var IsLevelNormal = glog.IsLevelNormal
var IsLevelSilent = glog.IsLevelSilent
var IsLevelVerbose = glog.IsLevelVerbose
var Log = glog.Log
var Normal = glog.Normal
var Prompt = glog.Prompt
var Verbose = glog.Verbose
