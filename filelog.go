// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"os"
	"fmt"
	"time"
)

var DEFUALT_HOUR_SUFFIX string = "20060102.15"
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

	// Rotate at linecount
	maxlines          int
	maxlines_curlines int

	// Rotate at size
	maxsize         int
	maxsize_cursize int

	// Rotate daily
	daily          bool
	daily_opendate int

	// Keep old logfiles (.001, .002, etc)
	rotate bool

	//新增按小时
	hourly         bool
	hourfilesuffix string
	hourinterval   int
	// after hourinterval hour unix-time sec
	nxtRotateSec int64
	curFileName  string
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
func NewFileLogWriter(fname string, rotate bool, hourly bool, hourFileSuffix string, hourInterval int) *FileLogWriter {
	w := &FileLogWriter{
		rec:      make(chan *LogRecord, LogBufferLength),
		rot:      make(chan bool),
		filename: fname,
		format:   "[%D %T] [%L] (%S) %M",
		rotate:   rotate,
		//新增小时
		hourly:         hourly,
		hourfilesuffix: hourFileSuffix,
		hourinterval:   hourInterval,
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
				n2 := time.Now()
				//if (w.maxlines > 0 && w.maxlines_curlines >= w.maxlines) ||
				//	(w.maxsize > 0 && w.maxsize_cursize >= w.maxsize) ||
				//	(w.daily && now.Day() != w.daily_opendate) {
				//	if err := w.intRotate(); err != nil {
				//		fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
				//		return
				//	}
				//}
				//新增小时
				if w.rotate && w.hourly && n2.Unix() >= w.nxtRotateSec {
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

	if w.rotate {
		//可以滚动
		//新增小时

		//是否按小时滚动
		if w.hourly {
			//可以小时滚动
			now := time.Now()
			format := w.hourfilesuffix
			if format == "" {
				format = DEFUALT_HOUR_SUFFIX
				w.hourfilesuffix = DEFUALT_HOUR_SUFFIX
			}
			//修改要写入的文件名称
			toWriteFile = w.filename + "." + fmt.Sprintf("%s", now.Format(format))
			//设置下一次滚动的时间
			interval := w.hourinterval
			if interval < 0 {
				w.hourinterval, interval = 1, 1
			}
			aft1hour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+interval, 0, 0, 0, time.Local)
			//aft1hour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second() + 5, 0, time.Local)
			w.nxtRotateSec = aft1hour.Unix()
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

	// Set the daily open date to the current date
	w.daily_opendate = now.Day()

	// initialize rotation values
	w.maxlines_curlines = 0
	w.maxsize_cursize = 0

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
func (w *FileLogWriter) SetHeadFoot(head, foot string) *FileLogWriter {
	w.header, w.trailer = head, foot
	if w.maxlines_curlines == 0 {
		fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: time.Now()}))
	}
	return w
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *FileLogWriter) SetRotateLines(maxlines int) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateLines: %v\n", maxlines)
	w.maxlines = maxlines
	return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *FileLogWriter) SetRotateSize(maxsize int) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
	w.maxsize = maxsize
	return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *FileLogWriter) SetRotateDaily(daily bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateDaily: %v\n", daily)
	w.daily = daily
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *FileLogWriter) SetRotate(rotate bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotate: %v\n", rotate)
	w.rotate = rotate
	return w
}

//新增小时
func (w *FileLogWriter) SetHourly(hourly bool) *FileLogWriter {
	w.hourly = hourly
	return w
}
func (w *FileLogWriter) SetHourFileSuffix(hourfilesuffix string) *FileLogWriter {
	w.hourfilesuffix = hourfilesuffix
	return w
}
func (w *FileLogWriter) SetHourInterval(hourinterval int) *FileLogWriter {
	w.hourinterval = hourinterval
	return w
}

// NewXMLLogWriter is a utility method for creating a FileLogWriter set up to
// output XML record log messages instead of line-based ones.
func NewXMLLogWriter(fname string, rotate bool) *FileLogWriter {
	return NewFileLogWriter(fname, rotate, false, "", 1).SetFormat(
		`	<record level="%L">
		<timestamp>%D %T</timestamp>
		<source>%S</source>
		<message>%M</message>
	</record>`).SetHeadFoot("<log created=\"%D %T\">", "</log>")
}
