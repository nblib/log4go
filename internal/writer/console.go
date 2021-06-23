package writer

import (
	"fmt"
	"github.com/nblib/log4go/internal/core"
	"os"
	"sync"
)

type Console chan *core.Record

var (
	doOnceDefaultConsole sync.Once
	DefaultConsoleWriter Console
)

func (c Console) run() {
	var timestr string
	var timestrAt int64

	stdOut := os.Stdout
	stdErr := os.Stderr

	for rec := range c {
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format("01/02/06 15:04:05"), at
		}
		if rec.Level <= core.INFO {
			_, _ = fmt.Fprint(stdOut, "[", timestr, "] [", core.LevelStrings[rec.Level], "] ", rec.Message, "\n")
		}
		if rec.Level >= core.WARN {
			_, _ = fmt.Fprint(stdErr, "[", timestr, "] [", core.LevelStrings[rec.Level], "] ", rec.Message, "\n")
		}
	}
}

func (c Console) Write(rec *core.Record) {
	c <- rec
}

func (c Console) Close() {
	close(c)
}

func NewConsoleWriter() (Writer, error) {
	doOnceDefaultConsole.Do(func() {
		c := make(Console, 1)
		go c.run()
		DefaultConsoleWriter = c
	})

	return DefaultConsoleWriter, nil
}
