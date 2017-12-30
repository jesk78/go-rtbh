package config

import (
	"fmt"

	"github.com/r3boot/go-rtbh/lib/logger"
)

const MYNAME string = "Config"

var log *logger.Logger

func NewConfig(l *logger.Logger, cfgfile string) (*Config, error) {
	log = l

	cfg := &Config{}

	if err := cfg.LoadFrom(cfgfile); err != nil {
		return nil, fmt.Errorf("NewConfig: %v", err)
	}

	cfg.CheckAndSetDefaults()

	if err := cfg.CompileRuleset(); err != nil {
		return nil, fmt.Errorf("NewConfig: %v", err)
	}

	return cfg, nil
}
