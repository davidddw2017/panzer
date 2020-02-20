package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type Config struct {
	Server   `ini:"server"`
	DBConfig `ini:"db"`
}

type Server struct {
	Host          string `ini:"host"`
	Port          int    `ini:"port"`
	ViewPattern   string `ini:"viewPattern"`
	StaticPattern string `ini:"staticPattern"`
	Env           string `ini:"env"`
}

type DBConfig struct {
	Host             string `ini:"host"`
	Port             int    `ini:"port"`
	Dbname           string `ini:"dbname"`
	User             string `ini:"user"`
	Passwd           string `ini:"passwd"`
	Charset          string `ini:"charset"`
	MaxOpenConns     int    `ini:"maxOpenConns"`
	MaxIdleConns     int    `ini:"maxIdleConns"`
	MaxLifetimeConns int    `ini:"maxLifetimeConns"`
}

// SystemConfig global config
var SystemConfig *Config

func init() {
	var err error
	SystemConfig, err = Load("config.ini")
	if err != nil {
		log.Println(err)
	}
}

// Load get config from ini
func Load(configFileMame string) (*Config, error) {
	cfg, err := ini.Load(configFileMame)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = cfg.MapTo(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
