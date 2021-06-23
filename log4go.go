package log4go

import "github.com/nblib/log4go/internal"

var ConfiguredLoggers map[string]Logger

type Logger interface {
	D(arg0 interface{}, args ...interface{})
	I(arg0 interface{}, args ...interface{})
	W(arg0 interface{}, args ...interface{})
	E(arg0 interface{}, args ...interface{})
	Debug(arg0 interface{}, args ...interface{})
	Info(arg0 interface{}, args ...interface{})
	Warn(arg0 interface{}, args ...interface{})
	Error(arg0 interface{}, args ...interface{})
}

type loggerWrapper struct {
	logger internal.LoggerWriter
}

func (l *loggerWrapper) D(arg0 interface{}, args ...interface{}) {

}

func (l *loggerWrapper) I(arg0 interface{}, args ...interface{}) {
	panic("implement me")
}

func (l *loggerWrapper) W(arg0 interface{}, args ...interface{}) {
	panic("implement me")
}

func (l *loggerWrapper) E(arg0 interface{}, args ...interface{}) {
	panic("implement me")
}

func (l *loggerWrapper) Debug(arg0 interface{}, args ...interface{}) {
	panic("implement me")
}

func (l *loggerWrapper) Info(arg0 interface{}, args ...interface{}) {
	panic("implement me")
}

func (l *loggerWrapper) Warn(arg0 interface{}, args ...interface{}) {
	panic("implement me")
}

func (l *loggerWrapper) Error(arg0 interface{}, args ...interface{}) {
	panic("implement me")
}
