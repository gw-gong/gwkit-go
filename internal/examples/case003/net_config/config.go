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

func (c *Config) UnmarshalConfig() error {
	c.BaseConfig.Mu.Lock()
	defer c.BaseConfig.Mu.Unlock()

	if err := c.BaseConfig.Viper.Unmarshal(&c); err != nil {
		return fmt.Errorf("unmarshal config failed: %w", err)
	}
	return nil
}

func (c *Config) ReloadConfig() {
	if err := c.UnmarshalConfig(); err != nil {
		log.Error("Failed to reload config", log.Err(err))
		return
	}

	log.Info("consul Config reloaded successfully", log.Any("config", NetCfg))
}
