package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

const REDIS_D_ADDR string = "localhost:6379"

const AMQP_D_ADDR string = "localhost:5672"
const AMQP_D_USER string = "go-rtbh"
const AMQP_D_PASS string = "go-rtbh"
const AMQP_D_EXCHANGE string = "amqp-input"

type Config struct {
	Amqp struct {
		Address  string `yaml:"address"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Exchange string `yaml:"exchange"`
	} `yaml:"amqp"`
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Database int64  `yaml:"database"`
	} `yaml:"redis"`
	Whitelist []struct {
		Address     string `yaml:"address"`
		Description string `yaml:"description"`
	} `yaml:"whitelist"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = &Config{}

	return
}

func (cfg *Config) CheckAndSetDefaults() {
	if cfg.Redis.Address == "" {
		cfg.Redis.Address = REDIS_D_ADDR
	}
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

	// Set defaults (if not set)
	cfg.CheckAndSetDefaults()

	return
}
