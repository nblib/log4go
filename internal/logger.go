package internal

import (
	"github.com/nblib/log4go/internal/core"
	"github.com/nblib/log4go/internal/writer"
)

type LoggerWriter interface {
	LogRecord(rec *core.Record)
}

type DefaultLogger struct {
	Level     core.LEVEL
	Forbidden map[string]struct{}
	Writers   []writer.Writer
}
type NormalLogger struct {
	Level   core.LEVEL
	Writers []writer.Writer
}

func NewDefaultLogger(level core.LEVEL, forb []string, ws []writer.Writer) *DefaultLogger {
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
		writers[0], _ = writer.NewConsoleWriter()
	}
	return &DefaultLogger{
		Level:     level,
		Forbidden: forbMap,
		Writers:   writers,
	}
}
func NewNormalLogger(level core.LEVEL, ws []writer.Writer) *NormalLogger {
	//writer
	var writers []writer.Writer
	if ws != nil && len(ws) > 0 {
		writers = make([]writer.Writer, 0, len(ws))
		for _, item := range ws {
			writers = append(writers, item)
		}
	} else {
		writers = make([]writer.Writer, 1)
		writers[0], _ = writer.NewConsoleWriter()
	}
	return &NormalLogger{
		Level:   level,
		Writers: writers,
	}
}
