package internal

import (
	"fmt"
	"github.com/nblib/log4go/internal/core"
	"github.com/nblib/log4go/internal/writer"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
)

const (
	DefaultLevel = core.INFO
)

//func LoadNormalLogger(cfg *ini.File) ([]*NormalLogger, error) {
//
//}

func LoadDefaultLogger(cfg *ini.File) (*DefaultLogger, error) {
	sectionName := "default"
	if cfg == nil {
		log.Fatalln("LoadDefaultSection args must not nil")
	}
	var (
		finalLevel  core.LEVEL
		finalForb   []string
		finalWriter []writer.Writer
	)

	//section
	section, err := cfg.GetSection(sectionName)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "log4go: not found default setting,will use [INFO] level, and send to [os.Stdout] and [os.Stderr]")
		//use default
		return NewDefaultLogger(DefaultLevel, nil, nil), nil
	}
	//level
	if level, err := loadRequireLevel(section); err != nil {
		return nil, err
	} else {
		finalLevel = level
	}
	//forbidden
	if key, err := section.GetKey("forbidden"); err != nil {
		//pass
	} else {
		forbidden := key.Strings(",")
		finalForb = forbidden
	}
	//writers
	childSections := section.ChildSections()
	if childSections == nil {
		finalWriter = nil
	} else {
		finalWriter = make([]writer.Writer, 0, len(childSections))
		for _, child := range childSections {
			w, err := parseWriter(child, sectionName)
			if err != nil {
				return nil, err
			}
			finalWriter = append(finalWriter, w)
		}
	}

	return NewDefaultLogger(finalLevel, finalForb, finalWriter), nil
}

func parseWriter(section *ini.Section, parentName string) (writer.Writer, error) {
	name := section.Name()
	name = strings.TrimLeft(name, parentName)
	name = strings.ToLower(name)
	switch name {
	case ".console":
		return writer.NewConsoleWriter()
	default:
		return nil, fmt.Errorf("log4go: unknown writer type: section:[%s],expect:[%s] got: [%s]", section.Name(), ".file,.console,.socket,...", name)
	}
}

func loadLevel(section *ini.Section, defaultVal core.LEVEL) core.LEVEL {
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
func loadRequireLevel(section *ini.Section) (core.LEVEL, error) {
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
