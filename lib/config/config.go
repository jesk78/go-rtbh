package config

import (
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "Config"

var Log logger.Log

func Setup(l logger.Log) (err error) {
	Log = l

	return
}

func New(cfgfile string) *Config {
	var cfg *Config
	var err error

	cfg = &Config{}

	if err = cfg.LoadFrom(cfgfile); err != nil {
		Log.Fatal(err)
	}

	cfg.CheckAndSetDefaults()

	if err = cfg.CompileRuleset(); err != nil {
		Log.Fatal(err)
	}

	return cfg
}
