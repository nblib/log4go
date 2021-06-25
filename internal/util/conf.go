package util

import (
	"fmt"
	"github.com/nblib/log4go/internal/core"
	"gopkg.in/ini.v1"
	"strings"
)

func LoadBoolConf(section *ini.Section, name string, def bool) bool {
	if section == nil {
		return def
	}
	if key, err := section.GetKey(name); err != nil {
		return def
	} else {
		return key.MustBool(def)
	}
}

func LoadLevel(section *ini.Section, defaultVal core.LEVEL) core.LEVEL {
	fieldName := "level"
	key, err := section.GetKey(fieldName)
	if err != nil {
		return defaultVal
	}
	loadLevel := key.String()
	loadLevel = strings.ToLower(loadLevel)
	level, ok := core.LevelMap[loadLevel]
	if !ok {
		return defaultVal
	}
	return level
}
func LoadRequireLevel(section *ini.Section) (core.LEVEL, error) {
	key, err := section.GetKey("level")
	if err != nil {
		return -1, cfgRequireErr(section.Name(), "level")
	}
	loadLevel := key.String()
	loadLevel = strings.ToLower(loadLevel)
	level, ok := core.LevelMap[loadLevel]
	if !ok {
		return -1, cfgUnexpectErr(section.Name(), "debug,info,warn,error", loadLevel)
	}
	return level, nil
}
func cfgRequireErr(name, require string) error {
	return fmt.Errorf("log4go: load cfg error: section:[%s], require:[%s]", name, require)
}
func cfgUnexpectErr(name, require, got string) error {
	return fmt.Errorf("log4go: unexpect cfg: section:[%s], require:[%s], got:[%s]", name, require, got)
}
