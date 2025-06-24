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
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	API struct {
		Key     string        `yaml:"key"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"api"`
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
