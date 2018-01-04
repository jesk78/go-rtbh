package config

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

var Ruleset []*regexp.Regexp

func (cfg *Config) CheckAndSetDefaults() {
}

func (cfg *Config) CompileRuleset() error {
	for _, ruleData := range cfg.Ruleset {
		re, err := regexp.Compile(ruleData)
		if err != nil {
			return fmt.Errorf("Config.CompileRuleset regexp.Compile: %v", err)
		}

		Ruleset = append(Ruleset, re)
	}

	return nil
}

func (cfg *Config) LoadFrom(fname string) error {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Config.LoadFrom ioutil.ReadFile: %v", err)
	}

	// Try to parse the data structure into yaml
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("Config.LoadFrom yaml.Unmarshal: %v", err)
	}

	return nil
}
