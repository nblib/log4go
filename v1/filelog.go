// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package v1

import (
	"fmt"
	"os"
	"time"
)

//The default suffix of the file saved by time
const DEFUALT_ROT_TIME_SUFFIX_HOUR string = "20060102_15"
const DEFUALT_ROT_TIME_SUFFIX_DAY string = "200601_02"

//rotation type
const ROT_TYPE_TIME_HOUR string = "hour"
const ROT_TYPE_TIME_DAY string = "day"

// This log writer sends output to a file
type FileLogWriter struct {
	rec chan *LogRecord
	rot chan bool

	// The opened file
	filename string
	file     *os.File

	// The logging format
	format string

	// File header/trailer
	header, trailer string

	// Keep old logfiles (.001, .002, etc)
	rotatable   bool
	rottype     string
	rotsuffix   string
	rotinterval int
	// get unix-time according to rotinterval and rottype
	nxtRotateTime int64
	curFileName   string
}

// This is the FileLogWriter's output method
func (w *FileLogWriter) LogWrite(rec *LogRecord) {
	w.rec <- rec
}

func (w *FileLogWriter) Close() {
	close(w.rec)
}

// NewFileLogWriter creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewFileLogWriter(fname string, rotatable bool, rottype string, rotsuffix string, rotinterval int) *FileLogWriter {
	w := &FileLogWriter{
		rec:      make(chan *LogRecord, LogBufferLength),
		rot:      make(chan bool),
		filename: fname,
		format:   "[%D %T] [%L] (%S) %M",

		rotatable:   rotatable,
		rottype:     rottype,
		rotsuffix:   rotsuffix,
		rotinterval: rotinterval,
	}

	// open the file for the first time
	if err := w.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
		return nil
	}

	go func() {
		defer func() {
			if w.file != nil {
				fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
				w.file.Close()
			}
		}()

		for {
			select {
			case rec, ok := <-w.rec:
				if !ok {
					return
				}

				if w.canRotate() {
					// 更换新文件
					if err := w.intRotate(); err != nil {
						fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
						return
					}
				}
				// 记录日志
				_, err := fmt.Fprint(w.file, FormatLogRecord(w.format, rec))
				if err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.curFileName, err)
					return
				}

				// Update the counts
				//w.maxlines_curlines++
				//w.maxsize_cursize += n
			}
		}
	}()

	return w
}

//是否到了转换的时候
func (w *FileLogWriter) canRotate() bool {
	if !w.rotatable {
		return false
	}
	if w.rottype == ROT_TYPE_TIME_DAY || w.rottype == ROT_TYPE_TIME_HOUR {
		now := time.Now()
		if now.Unix() >= w.nxtRotateTime {
			return true
		}
	}
	return false
}

// Request that the logs rotate
func (w *FileLogWriter) Rotate() {
	w.rot <- true
}

// If this is called in a threaded context, it MUST be synchronized
func (w *FileLogWriter) intRotate() error {
	// 开始更换文件,首先关闭当前记录的日志文件
	if w.file != nil {
		fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
		w.file.Close()
	}
	/* 然后判断是否可以滚动
	滚动当前分为三种:日期,文件大小,行数.
	日期以小时为单位,比如一天则为24,一个小时为1,
	文件大小以M为单位.
	行数以k为单位.*/

	//写入的文件名,开始为设置的文件名 filename
	toWriteFile := w.filename

	if w.rotatable {
		//可以滚动
		//新增小时

		//是否按小时滚动
		if w.rottype == ROT_TYPE_TIME_DAY || w.rottype == ROT_TYPE_TIME_HOUR {
			//可以天滚动
			now := time.Now()
			if w.rotsuffix == "" {
				if w.rottype == ROT_TYPE_TIME_DAY {
					w.rotsuffix = DEFUALT_ROT_TIME_SUFFIX_DAY
				} else if w.rottype == ROT_TYPE_TIME_HOUR {
					w.rotsuffix = DEFUALT_ROT_TIME_SUFFIX_HOUR
				}
			}
			suffix := w.rotsuffix
			//修改要写入的文件名称
			toWriteFile = w.filename + fmt.Sprintf("%s", now.Format(suffix))
			//设置下一次滚动的时间
			interval := w.rotinterval
			if interval <= 0 {
				w.rotinterval, interval = 1, 1
			}
			var aftTime time.Time
			if w.rottype == ROT_TYPE_TIME_DAY {
				aftTime = time.Date(now.Year(), now.Month(), now.Day()+interval, 0, 0, 0, 0, time.Local)
			} else if w.rottype == ROT_TYPE_TIME_HOUR {
				aftTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+interval, 0, 0, 0, time.Local)
			}
			w.nxtRotateTime = aftTime.Unix()
		}

	}
	// Open the log file
	fd, err := os.OpenFile(toWriteFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	w.curFileName = toWriteFile
	w.file = fd

	now := time.Now()
	fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: now}))

	return nil
}

// Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w *FileLogWriter) SetFormat(format string) *FileLogWriter {
	w.format = format
	return w
}

// Set the logfile header and footer (chainable).  Must be called before the first log
// message is written.  These are formatted similar to the FormatLogRecord (e.g.
// you can use %D and %T in your header/footer for date and time).
//func (w *FileLogWriter) SetHeadFoot(head, foot string) *FileLogWriter {
//	w.header, w.trailer = head, foot
//	if w.maxlines_curlines == 0 {
//		fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: time.Now()}))
//	}
//	return w
//}
//
//// Set rotate at linecount (chainable). Must be called before the first log
//// message is written.
//func (w *FileLogWriter) SetRotateLines(maxlines int) *FileLogWriter {
//	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateLines: %v\n", maxlines)
//	w.maxlines = maxlines
//	return w
//}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
//func (w *FileLogWriter) SetRotateSize(maxsize int) *FileLogWriter {
//	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
//	w.maxsize = maxsize
//	return w
//}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *FileLogWriter) SetRotatable(rotatable bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotate: %v\n", rotate)
	w.rotatable = rotatable
	return w
}

func (w *FileLogWriter) SetRotSuffix(rotsuffix string) *FileLogWriter {
	w.rotsuffix = rotsuffix
	return w
}
func (w *FileLogWriter) SetRotInterval(rotinterval int) *FileLogWriter {
	w.rotinterval = rotinterval
	return w
}

// NewXMLLogWriter is a utility method for creating a FileLogWriter set up to
// output XML record log messages instead of line-based ones.
func NewXMLLogWriter(fname string, rotate bool) *FileLogWriter {
	return NewFileLogWriter(fname, rotate, "day", "", 1).SetFormat(
		`	<record level="%L">
		<timestamp>%D %T</timestamp>
		<source>%S</source>
		<message>%M</message>
	</record>`)
}
