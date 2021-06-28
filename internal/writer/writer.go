package writer

import (
	"github.com/nblib/log4go/v2/internal/core"
)

type Writer interface {
	Write(rec *core.Record)
	Close()
}
type WriterRoot struct {
	OutTime    bool
	OutSource  bool
	OutLogName bool
	Level      core.LEVEL
}
