package net_config

import (
	"fmt"
	"sync"
	"time"

	"github.com/gw-gong/gwkit-go/hot_cfg"
	"github.com/gw-gong/gwkit-go/log"
)

var (
	NetCfg *Config
	once   sync.Once
)

type Config struct {
	*hot_cfg.BaseConfig
	Database struct {
		Host     string `yaml:"host" mapstructure:"host"`
		Port     int    `yaml:"port" mapstructure:"port"`
		Username string `yaml:"username" mapstructure:"username"`
		Password string `yaml:"password" mapstructure:"password"`
	} `yaml:"database" mapstructure:"database"`
	API struct {
		Key     string        `yaml:"key" mapstructure:"key"`
		Timeout time.Duration `yaml:"timeout" mapstructure:"timeout"`
	} `yaml:"api" mapstructure:"api"`
}

func InitNetConfig(consulAddr, consulKey, configType string, reloadTime int) error {
	var err error
	once.Do(func() {
		NetCfg = &Config{}
		NetCfg.BaseConfig, err = hot_cfg.NewBaseConfig(
			hot_cfg.WithConsulConfig(consulAddr, consulKey, configType, reloadTime),
		)
		if err != nil {
			err = fmt.Errorf("init base config failed: %w", err)
		}
	})
	return err
}

func (c *Config) LoadConfig() {
	c.BaseConfig.Mu.Lock()
	defer c.BaseConfig.Mu.Unlock()

	if err := c.BaseConfig.Viper.Unmarshal(&c); err != nil {
		log.Error("unmarshal config failed", log.Err(err))
		return
	}

	log.Info("LoadConfig", log.Any("config", c))
}
