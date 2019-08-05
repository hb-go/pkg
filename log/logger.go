package log

var logger Logger

func init() {
	logger = NewLogger()
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

func SetLogger(l Logger) {
	logger = l
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	logger.Info(v...)
}
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	logger.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	logger.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	logger.Panic(v...)
}
func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}
