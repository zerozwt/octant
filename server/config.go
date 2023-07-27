package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zerozwt/swe"
	"gopkg.in/yaml.v3"
)

type ConfigService struct {
	Core      bool `yaml:"core"`
	Collector bool `yaml:"collector"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type Config struct {
	LocalHost bool          `yaml:"localhost"`
	Port      uint16        `yaml:"port"`
	DbEngine  string        `yaml:"db_engine"`
	MySQL     string        `yaml:"mysql"`
	SQLite    string        `yaml:"sqlite"`
	WebDir    string        `yaml:"www_dir"`
	Service   ConfigService `yaml:"service"`
	Etcd      []string      `yaml:"etcd"`
	Log       LogConfig     `yaml:"log"`
}

func (c Config) IsMySQL() bool  { return c.DbEngine == "mysql" }
func (c Config) IsSQLite() bool { return c.DbEngine == "sqlite" }
func (c Config) LogLevel() swe.LogLevel {
	switch c.Log.Level {
	case "debug":
		return swe.LOG_DEBUG
	case "info":
		return swe.LOG_INFO
	case "warn":
		return swe.LOG_WARN
	case "error":
		return swe.LOG_ERROR
	}
	return swe.LOG_INFO
}

var gConfig Config = Config{
	LocalHost: true,
	Port:      6080,
	DbEngine:  "sqlite",
	SQLite:    "./octant.db",
	WebDir:    "./dist",
	Service: ConfigService{
		Core:      true,
		Collector: true,
	},
	Log: LogConfig{
		Level: "debug",
		File:  "",
	},
}

func LoadConfig(file string) error {
	confData, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(confData, &gConfig)
	if err != nil {
		return err
	}

	gConfig.DbEngine = strings.ToLower(gConfig.DbEngine)

	if gConfig.Port == 0 {
		return fmt.Errorf("invalid port %d", gConfig.Port)
	}

	if gConfig.DbEngine != "sqlite" && gConfig.DbEngine != "mysql" {
		return fmt.Errorf("invalid db engine: %s", gConfig.DbEngine)
	}

	if !(gConfig.Service.Core || gConfig.Service.Collector) {
		return fmt.Errorf("no service")
	}

	gConfig.Log.Level = strings.ToLower(gConfig.Log.Level)

	return nil
}
