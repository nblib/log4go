// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// This log writer sends output to a socket
type SocketLogWriter chan *LogRecord

var UDP_FAILED_WAITTIME time.Duration = 1
var TCP_RECONN_WAITTIME time.Duration = 3

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
		fmt.Fprintf(os.Stderr, "NewSocketLogWriter(Start connect): %s,err: %v\n", hostport, err)
		//return nil
	}

	w := SocketLogWriter(make(chan *LogRecord, LogBufferLength))

	go func() {
		isTcp := strings.EqualFold(proto, "tcp")
		defer func() {
			if sock != nil && isTcp {
				sock.Close()
			}
		}()
		if sock == nil {
			//重连
			sock = blockForReconnect(proto, hostport)
		}
		for rec := range w {
			// Marshall into yyyy-MM-dd HH:mm:ss message
			message := rec.toSTR()
			_, err = sock.Write(message)
			if err != nil {

				fmt.Fprintf(os.Stderr, "SocketLogWriter(At writing): %s,err: %v \n", hostport, err)
				//发送失败,等待一段时间重试或重连
				if isTcp {
					_ = sock.Close()
					sock = blockForReconnect(proto, hostport)
				} else {
					//不是tcp,等待一会再试
					time.Sleep(UDP_FAILED_WAITTIME * time.Second)
				}
			}
		}
	}()

	return w
}

/**
阻塞等待重连接
*/
func blockForReconnect(proto string, addr string) (sock net.Conn) {
	var err error
	for {
		sock, err = net.Dial(proto, addr)
		if err == nil {
			return
		} else {
			time.Sleep(TCP_RECONN_WAITTIME * time.Second)
		}
	}
}
