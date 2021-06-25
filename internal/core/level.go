package core

type LEVEL int8

const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
)

var LevelStrings = [...]string{"DEBUG", "INFO", "WARN", "ERROR"}
var LevelMap = map[string]LEVEL{
	"debug": DEBUG,
	"info":  INFO,
	"warn":  WARN,
	"error": ERROR,
}
