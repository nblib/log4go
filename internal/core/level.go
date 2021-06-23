package core

type LEVEL int8

const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
)

var LevelStrings = [...]string{"FNST", "FINE", "DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT"}
var LevelMap = map[string]LEVEL{
	"debug": DEBUG,
	"info":  INFO,
	"warn":  WARN,
	"error": ERROR,
}
