package writer

import (
	"fmt"
	"github.com/nblib/log4go/v2/internal/core"
	"github.com/nblib/log4go/v2/internal/util"
	"gopkg.in/ini.v1"
	"os"
	"strings"
	"sync"
)

type Console struct {
	WriterRoot
}

var (
	doOnceDefaultConsole sync.Once
	DefaultConsoleWriter Console
)

func (c Console) Write(rec *core.Record) {
	builder := &strings.Builder{}
	if c.OutLogName {
		builder.WriteString("<")
		builder.WriteString(rec.LoggerName)
		builder.WriteString(">")
		builder.WriteString(" ")
	}
	if c.OutSource && rec.Source != "" {
		builder.WriteString(rec.Source)
		builder.WriteString(" ")
	}
	if c.OutTime {
		timeStr := rec.Created.Format("2006-01-02 15:04:05.000")
		builder.WriteString(timeStr)
		builder.WriteString(" ")
	}
	builder.WriteString("[")
	builder.WriteString(core.LevelStrings[rec.Level])
	builder.WriteString("]")
	builder.WriteString(" ")
	builder.WriteString(rec.Message)
	builder.WriteString("\n")
	if rec.Level <= core.INFO {
		_, _ = fmt.Fprint(os.Stdout, builder.String())
	}
	if rec.Level >= core.WARN {
		_, _ = fmt.Fprint(os.Stderr, builder.String())
	}
}

func (c Console) Close() {
	//pass
}

func NewConsoleWriter(section *ini.Section) (Writer, error) {
	doOnceDefaultConsole.Do(func() {
		finalOutLogName := util.LoadBoolConf(section, "out_logname", true)
		finalOutSource := util.LoadBoolConf(section, "out_source", true)
		finalOutTime := util.LoadBoolConf(section, "out_time", true)

		c := Console{
			WriterRoot: WriterRoot{
				OutTime:    finalOutTime,
				OutSource:  finalOutSource,
				OutLogName: finalOutLogName,
			},
		}
		DefaultConsoleWriter = c
	})

	return DefaultConsoleWriter, nil
}
