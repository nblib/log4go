// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"net"
	"os"
	"time"
)

// This log writer sends output to a socket
type SocketLogWriter chan *LogRecord

var UDP_FAILED_WAITTIME time.Duration = 3
// This is the SocketLogWriter's output method
func (w SocketLogWriter) LogWrite(rec *LogRecord) {
	//小于缓冲区,才发送,不小于,不发送
	if len(w) < LogBufferLength {
		w <- rec
	}
}

func (w SocketLogWriter) Close() {
	close(w)
}

func NewSocketLogWriter(proto, hostport string) SocketLogWriter {
	sock, err := net.Dial(proto, hostport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewSocketLogWriter(Start connect): %s\n", hostport, err)
		return nil
	}

	w := SocketLogWriter(make(chan *LogRecord, LogBufferLength))

	go func() {
		defer func() {
			if sock != nil && proto == "tcp" {
				sock.Close()
			}
		}()

		for rec := range w {
			// Marshall into yyyy-MM-dd HH:mm:ss message
			message := rec.toSTR()
			_, err = sock.Write(message)
			if err != nil {
				fmt.Fprint(os.Stderr, "SocketLogWriter(At writing): %s", hostport, err)
				//发送失败,等待一段时间重试
				time.Sleep(UDP_FAILED_WAITTIME * time.Second)
			}
		}
	}()

	return w
}
