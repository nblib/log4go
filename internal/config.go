package internal

import (
	"fmt"
	"github.com/nblib/log4go/v2/internal/core"
	"github.com/nblib/log4go/v2/internal/util"
	"github.com/nblib/log4go/v2/internal/writer"
	"gopkg.in/ini.v1"
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
		return NewDefaultLogger(DefaultLevel, nil, nil, true), nil
	}
	var (
		finalLevel        core.LEVEL
		finalForb         []string
		finalWriter       []writer.Writer
		finalEnableSource bool
	)

	//section
	section, err := cfg.GetSection(sectionName)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "log4go: not found default setting,will use [INFO] level, and send to [os.Stdout] and [os.Stderr]")
		//use default
		return NewDefaultLogger(finalLevel, finalForb, finalWriter, finalEnableSource), nil
	}
	//level
	if level, err := util.LoadRequireLevel(section); err != nil {
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
	//out condition
	finalEnableSource = util.LoadBoolConf(section, "enable_source", true)

	return NewDefaultLogger(finalLevel, finalForb, finalWriter, finalEnableSource), nil
}

func parseWriter(section *ini.Section, parentName string) (writer.Writer, error) {
	name := section.Name()
	name = strings.TrimLeft(name, parentName)
	name = strings.ToLower(name)
	switch name {
	case ".console":
		return writer.NewConsoleWriter(section)
	default:
		return nil, fmt.Errorf("log4go: unknown writer type: section:[%s],expect:[%s] got: [%s]", section.Name(), ".file,.console,.socket,...", name)
	}
}
