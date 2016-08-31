package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"regexp"
)

var Ruleset []*regexp.Regexp

func NewConfig() (cfg *Config, err error) {
	cfg = &Config{}

	return
}

func (cfg *Config) CheckAndSetDefaults() {
}

func (cfg *Config) CompileRuleset() (err error) {
	var re *regexp.Regexp

	for _, rule_data := range cfg.Ruleset {
		if re, err = regexp.Compile(rule_data); err != nil {
			err = errors.New(MYNAME + ": Invalid regexp pattern: " + err.Error())
			return
		}

		Ruleset = append(Ruleset, re)
	}

	return
}

func (cfg *Config) LoadFrom(fname string) (err error) {
	var fd *os.File
	var fs os.FileInfo
	var data []byte

	// Check if the file exists and save the stat() info
	if fs, err = os.Stat(fname); err != nil {
		err = errors.New("[config.LoadFrom]: os.Stat(" + fname + "): " + err.Error())
		return
	}

	// Create a buffer large enough to hold the size of the complete config
	data = make([]byte, fs.Size())

	// Open the file and defer close it
	if fd, err = os.Open(fname); err != nil {
		err = errors.New("[config.LoadFrom]: os.Open(" + fname + "): " + err.Error())
		return
	}
	defer fd.Close()

	// Read in the configuration file
	if _, err = fd.Read(data); err != nil {
		err = errors.New("[config.LoadFrom]: fd.Read(): " + err.Error())
		return
	}

	// Try to parse the data structure into yaml
	if err = yaml.Unmarshal(data, cfg); err != nil {
		err = errors.New("[config.LoadFrom]: yaml.Unmarshal(): " + err.Error())
		return
	}

	return
}
