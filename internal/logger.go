package internal

import (
	"github.com/nblib/log4go/internal/core"
	"github.com/nblib/log4go/internal/writer"
)

type LoggerWriter interface {
	LogRecord(rec *core.Record)
	Level() core.LEVEL
}

type DefaultLogger struct {
	Name         string
	level        core.LEVEL
	Forbidden    map[string]struct{}
	Writers      []writer.Writer
	EnableSource bool
}
type NormalLogger struct {
	Name    string
	Level   core.LEVEL
	Writers []writer.Writer
}

func (d DefaultLogger) LogRecord(rec *core.Record) {
	for _, w := range d.Writers {
		w.Write(rec)
	}
}

func (d DefaultLogger) Level() core.LEVEL {
	return d.level
}

func NewDefaultLogger(level core.LEVEL, forb []string, ws []writer.Writer, enableSource bool) *DefaultLogger {
	//forbidden
	var forbMap map[string]struct{}
	if forb != nil && len(forb) > 0 {
		forbMap = make(map[string]struct{}, len(forb))
		for _, item := range forb {
			forbMap[item] = struct{}{}
		}
	} else {
		forbMap = nil
	}
	//writer
	var writers []writer.Writer
	if ws != nil && len(ws) > 0 {
		writers = ws
	} else {
		writers = make([]writer.Writer, 1)
		writers[0], _ = writer.NewConsoleWriter(nil)
	}
	return &DefaultLogger{
		Name:         "default",
		level:        level,
		Forbidden:    forbMap,
		Writers:      writers,
		EnableSource: enableSource,
	}
}
func NewNormalLogger(name string, level core.LEVEL, ws []writer.Writer) *NormalLogger {
	//writer
	var writers []writer.Writer
	if ws != nil && len(ws) > 0 {
		writers = make([]writer.Writer, 0, len(ws))
		for _, item := range ws {
			writers = append(writers, item)
		}
	} else {
		writers = make([]writer.Writer, 1)
		writers[0], _ = writer.NewConsoleWriter(nil)
	}
	return &NormalLogger{
		Name:    name,
		Level:   level,
		Writers: writers,
	}
}
