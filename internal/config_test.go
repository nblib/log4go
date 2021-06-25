package internal

import (
	"github.com/nblib/log4go/internal/core"
	"gopkg.in/ini.v1"
	"testing"
)

func TestLoadDefaultLogger(t *testing.T) {
	t.Run("no_default", func(t *testing.T) {
		cfgData := []byte(`
`)
		cfg, err := ini.Load(cfgData)
		if err != nil {
			t.Fatal(err)
		}
		logger, err := LoadDefaultLogger(cfg)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(logger)
	})
	t.Run("only_level", func(t *testing.T) {
		runFunc := func(cfgData []byte, expected core.LEVEL, requireErr bool) {
			cfg, err := ini.Load(cfgData)
			if err != nil {
				t.Fatal(err)
			}
			logger, err := LoadDefaultLogger(cfg)
			if err != nil && !requireErr {
				t.Fatal(err)
			}
			if requireErr {
				if err != nil {
					t.Logf("expect err: %v\n", err)
				}
				if err == nil {
					t.Fatal("require error, but no error")
				}
			} else {
				if logger == nil || logger.Level() != expected {
					t.Fatal("load level error")
				}
			}
		}
		runFunc([]byte("[default]\nlevel=ERROR"), core.ERROR, false)
		runFunc([]byte("[default]\nlevel=error"), core.ERROR, false)
		runFunc([]byte("[default]\nlevel=info"), core.INFO, false)
		runFunc([]byte("[default]\nlevel=debug"), core.DEBUG, false)
		runFunc([]byte("[default]\nlevel=warn"), core.WARN, false)
		runFunc([]byte("[default]\nlevel=w"), core.WARN, true)
		runFunc([]byte("[default]\nlevel="), core.WARN, true)
	})
	t.Run("forbidden", func(t *testing.T) {
		runFunc := func(cfgData []byte, expect []string) {
			cfg, err := ini.Load(cfgData)
			if err != nil {
				t.Fatal(err)
			}
			logger, err := LoadDefaultLogger(cfg)
			if err != nil {
				t.Fatal(err)
			}
			gotForb := logger.Forbidden
			if gotForb != nil {
				if expect == nil {
					t.Fatal("ford not nil, require nil")
				}
				for _, item := range expect {
					if _, ok := gotForb[item]; !ok {
						t.Fatal("ford not eq expect")
					}
				}
			} else {
				if expect != nil {
					t.Fatal("ford is nil, require not nil")
				}
			}
		}
		runFunc([]byte("[default]\nlevel=info\nforbidden=a,b,c"), []string{"a", "b", "c"}[:])
		runFunc([]byte("[default]\nlevel=info\nforbidden=a,b,"), []string{"a", "b"}[:])
		runFunc([]byte("[default]\nlevel=info\nforbidden="), nil)
	})
}
