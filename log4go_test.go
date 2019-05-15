// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"errors"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestSocket(t *testing.T) {
	fmt.Println("nihao")
	l4ginfo := make(Logger)
	l4ginfo.LoadConfiguration("conf/l4gvisitinfo.xml")

	for {
		l4ginfo.Error("abc 23")
		fmt.Println("send...")
		time.Sleep(1 * time.Second)
	}

	time.Sleep(5 * time.Second)
}

func TestFprint(t *testing.T) {
	fmt.Fprintf(os.Stdout, "this is %s, err: %v \n", "1.22", errors.New("nihao"))
}

func TestSocketServer(t *testing.T) {
	listener, err := net.Listen("tcp", "10.169.43.52:3577")
	if err != nil {
		panic(err)
	}
	for {
		_, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("connect comming")

	}

}
