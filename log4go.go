package log4go

import (
	"fmt"
	"github.com/nblib/log4go/internal"
	"github.com/nblib/log4go/internal/core"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	ConfiguredLoggers map[string]Logger
	DefaultLogger     *loggerWrapper
)

func init() {
	initLog4go()
}
func initLog4go() {
	ConfiguredLoggers = make(map[string]Logger)

	configPath := "./log4go.ini"
	if getwd, err := os.Getwd(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "get work dir error, use path: ", configPath)
	} else {
		configPath = path.Join(getwd, configPath)
	}
	cfgFile, err := ini.Load(configPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "log4go: Not found config file from [%s], use default config\n", configPath)
		_ = wrapDefaultLogger(nil)
	} else {
		if err := wrapDefaultLogger(cfgFile); err != nil {
			log.Fatalf("log4go: load error %v\n", err)
		}
	}
}

func wrapDefaultLogger(cfgFile *ini.File) error {
	loggerWriter, err := internal.LoadDefaultLogger(cfgFile)
	if err != nil {
		return err
	}
	DefaultLogger = &loggerWrapper{loggerWriter: loggerWriter,
		Name: loggerWriter.Name,
	}
	ConfiguredLoggers[""] = DefaultLogger
	return nil
}

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
	loggerWriter internal.LoggerWriter
	recordPool   sync.Pool
	l            sync.RWMutex
	Name         string
}

func (l *loggerWrapper) D(arg0 interface{}, args ...interface{}) {
	l.log(core.DEBUG, arg0, args)
}

func (l *loggerWrapper) I(arg0 interface{}, args ...interface{}) {
	l.log(core.INFO, arg0, args)

}

func (l *loggerWrapper) W(arg0 interface{}, args ...interface{}) {
	l.log(core.WARN, arg0, args)

}

func (l *loggerWrapper) E(arg0 interface{}, args ...interface{}) {
	l.log(core.ERROR, arg0, args)

}

func (l *loggerWrapper) Debug(arg0 interface{}, args ...interface{}) {
	l.log(core.DEBUG, arg0, args)

}

func (l *loggerWrapper) Info(arg0 interface{}, args ...interface{}) {
	l.log(core.INFO, arg0, args)

}

func (l *loggerWrapper) Warn(arg0 interface{}, args ...interface{}) {
	l.log(core.WARN, arg0, args)

}

func (l *loggerWrapper) Error(arg0 interface{}, args ...interface{}) {
	l.log(core.ERROR, arg0, args)
}

func (l *loggerWrapper) log(level core.LEVEL, arg0 interface{}, args ...interface{}) {
	if !l.checkLevel(level) {
		DefaultLogger.log(level, arg0, args)
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		l.intLogf(level, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		l.intLogc(level, first)
	default:
		// Build a format string so that it will be similar to Sprint
		l.intLogf(level, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func (l *loggerWrapper) checkLevel(level core.LEVEL) bool {
	return l.loggerWriter.Level() <= level
}

/******* Logging *******/
// Send a formatted log message internally
func (l *loggerWrapper) intLogf(lvl core.LEVEL, format string, args ...interface{}) {

	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(3)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// Make the log record
	rec := &core.Record{
		Level:      lvl,
		Created:    time.Now(),
		Source:     src,
		Message:    msg,
		LoggerName: l.Name,
	}
	l.loggerWriter.LogRecord(rec)
}

// Send a closure log message internally
func (l *loggerWrapper) intLogc(lvl core.LEVEL, closure func() string) {
	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(3)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	// Make the log record
	rec := &core.Record{
		Level:      lvl,
		Created:    time.Now(),
		Source:     src,
		Message:    closure(),
		LoggerName: l.Name,
	}
	l.loggerWriter.LogRecord(rec)
}

func D(arg0 interface{}, args ...interface{}) {
	DefaultLogger.log(core.DEBUG, arg0, args...)
}

func I(arg0 interface{}, args ...interface{}) {
	DefaultLogger.log(core.INFO, arg0, args...)

}

func W(arg0 interface{}, args ...interface{}) {
	DefaultLogger.log(core.WARN, arg0, args...)

}

func E(arg0 interface{}, args ...interface{}) {
	DefaultLogger.log(core.ERROR, arg0, args...)

}

func Debug(arg0 interface{}, args ...interface{}) {
	DefaultLogger.log(core.DEBUG, arg0, args...)

}

func Info(arg0 interface{}, args ...interface{}) {
	DefaultLogger.log(core.INFO, arg0, args...)

}

func Warn(arg0 interface{}, args ...interface{}) {
	DefaultLogger.log(core.WARN, arg0, args...)

}

func Error(arg0 interface{}, args ...interface{}) {
	DefaultLogger.log(core.ERROR, arg0, args...)
}
