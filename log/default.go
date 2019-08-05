package log

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

const (
	DEBUG Lvl = iota
	INFO
	WARN
	ERROR
	OFF
	fatalLvl
	panicLvl
)

type (
	Lvl uint
)

func (lvl Lvl) ColorString() string {
	switch lvl {
	case DEBUG:
		return color.WhiteString("DEBUG")
	case INFO:
		return color.GreenString("INFO")
	case WARN:
		return color.YellowString("WARN")
	case ERROR:
		return color.RedString("ERROR")
	case fatalLvl:
		return color.HiRedString("FATAL")
	case panicLvl:
		return color.HiRedString("PANIC")
	default:
		return color.WhiteString("-")
	}
}

func (lvl Lvl) String() string {
	switch lvl {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case fatalLvl:
		return "FATAL"
	case panicLvl:
		return "PANIC"
	default:
		return "-"
	}
}

type defaultLogger struct {
	*log.Logger
	level       Lvl
	colorEnable bool
	calldepth   int
}

func (l *defaultLogger) SetLevel(lvl Lvl) {
	l.level = lvl
}

func (l *defaultLogger) SetColor(enable bool) {
	l.colorEnable = enable
}

func (l *defaultLogger) SetCalldepth(calldepth int) {
	l.calldepth = calldepth
}

func (l *defaultLogger) Debug(v ...interface{}) {
	l.output(DEBUG, v...)
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	l.outputf(DEBUG, format, v...)
}

func (l *defaultLogger) Info(v ...interface{}) {
	l.output(INFO, v...)
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	l.outputf(INFO, format, v...)
}

func (l *defaultLogger) Warn(v ...interface{}) {
	l.output(WARN, v...)
}

func (l *defaultLogger) Warnf(format string, v ...interface{}) {
	l.outputf(WARN, format, v...)
}

func (l *defaultLogger) Error(v ...interface{}) {
	l.output(ERROR, v...)
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	l.outputf(ERROR, format, v...)
}

func (l *defaultLogger) Fatal(v ...interface{}) {
	l.output(fatalLvl, v...)
	os.Exit(1)
}

func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	l.outputf(fatalLvl, format, v...)
	os.Exit(1)
}

func (l *defaultLogger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.output(panicLvl, s)
	panic(s)
}

func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.output(panicLvl, s)
	panic(s)
}

func (l *defaultLogger) output(lvl Lvl, v ...interface{}) {
	if lvl < l.level {
		return
	}
	l.Output(l.calldepth, l.format(lvl, fmt.Sprint(v...)))
}

func (l *defaultLogger) outputf(lvl Lvl, format string, v ...interface{}) {
	if lvl < l.level {
		return
	}
	l.Output(l.calldepth, l.format(lvl, fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) format(lvl Lvl, msg string) string {
	if l.colorEnable {
		return fmt.Sprintf("[%s] %s", lvl.ColorString(), msg)
	}

	return fmt.Sprintf("[%s] %s", lvl.String(), msg)
}

func NewLogger() *defaultLogger {
	return &defaultLogger{
		Logger:      log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
		level:       INFO,
		calldepth:   4,
		colorEnable: true,
	}
}
