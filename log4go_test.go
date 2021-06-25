package log4go

import (
	"errors"
	"fmt"
	"testing"
)

func TestInitLog4go(t *testing.T) {
	I("this is test Info")
	I("this is test Info")
	I("this is test Info")
	I("this is test Info")
	I("this is test Info")
	E("this is test Err %v", errors.New("test errr"))
}
func TestFmt(t *testing.T) {
	sprint := fmt.Sprint(errors.New("test errr"))
	fmt.Println(sprint)
	t.Fatal(errors.New("test err"))
}
