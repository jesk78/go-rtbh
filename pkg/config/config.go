package config

import (
	"fmt"

	"github.com/r3boot/go-rtbh/pkg/logger"
)

const MYNAME string = "Config"

var log *logger.Logger

func New(l *logger.Logger, cfgfile string) (*Config, error) {
	log = l

	cfg := &Config{}

	if err := cfg.LoadFrom(cfgfile); err != nil {
		return nil, fmt.Errorf("config.New: %v", err)
	}

	cfg.CheckAndSetDefaults()

	if err := cfg.CompileRuleset(); err != nil {
		return nil, fmt.Errorf("config.New: %v", err)
	}

	return cfg, nil
}
