package core

import (
	"bytes"
	"strconv"
	"time"
)

/****** LogRecord ******/

// A LogRecord contains all of the pertinent information for each message
type Record struct {
	Level      LEVEL     // The log level
	Created    time.Time // The time at which the log message was created (nanoseconds)
	Source     string    // The message source
	Message    string    // The log message
	LoggerName string    //The logger name
}

func (rec *Record) toJson() []byte {
	var buf bytes.Buffer
	buf.WriteString("{\"Level\":")
	buf.WriteString(strconv.Itoa(int(rec.Level)))
	buf.WriteString(",")

	buf.WriteString("\"Created\":")
	timeJson, _ := rec.Created.MarshalJSON()
	buf.Write(timeJson)
	buf.WriteString(",")

	buf.WriteString("\"LoggerName\":")
	buf.WriteString("\"")
	buf.WriteString(rec.LoggerName)
	buf.WriteString("\"")
	buf.WriteString(",")

	buf.WriteString("\"Source\":")
	buf.WriteString("\"")
	buf.WriteString(rec.Source)
	buf.WriteString("\"")
	buf.WriteString(",")
	buf.WriteString("\"Message\":")
	buf.WriteString("\"")
	buf.WriteString(rec.Message)
	buf.WriteString("\"")
	buf.WriteString("}")
	return buf.Bytes()
}
func (rec *Record) toSTR() []byte {
	var buf bytes.Buffer

	dateTime := rec.Created.Format("2006-01-02 15:04:05")
	buf.WriteString(dateTime)

	buf.WriteString(" ")

	buf.WriteString(rec.Message)
	buf.WriteString(" \n")
	return buf.Bytes()
}
