package zap

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/hb-go/pkg/log"
)

func ReplaceLogger(logger *zap.Logger) {
	zl := &zapLogger{logger: logger}
	log.SetLogger(zl)
}

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Debug(v ...interface{}) {
	l.logger.Debug(fmt.Sprint(v...))
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, v...))
}

func (l *zapLogger) Info(v ...interface{}) {
	l.logger.Info(fmt.Sprint(v...))
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l *zapLogger) Warn(v ...interface{}) {
	l.logger.Warn(fmt.Sprint(v...))
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, v...))
}

func (l *zapLogger) Error(v ...interface{}) {
	l.logger.Error(fmt.Sprint(v...))
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

func (l *zapLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(fmt.Sprint(v...))
}

func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, v...))
}

func (l *zapLogger) Panic(v ...interface{}) {
	l.logger.Panic(fmt.Sprint(v...))
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, v...))
}
