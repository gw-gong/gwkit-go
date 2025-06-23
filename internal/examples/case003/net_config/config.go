package net_config

import (
	"fmt"
	"sync"
	"time"

	"github.com/gw-gong/gwkit-go/hot_cfg"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"
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

func (c *Config) Watch() {
	go gwkit_common.WithRecover(func() {
		c.watchConfig()
	})
}

func (c *Config) watchConfig() {
	ticker := time.NewTicker(time.Duration(c.BaseConfig.ConsulConfig.ReloadTime) * time.Second)
	defer ticker.Stop()

	// Record initial configuration hash
	lastConfigHash := hot_cfg.CalculateConfigHash(c.BaseConfig.Viper)
	log.Info("Initial configuration hash", log.Any("lastConfigHash", lastConfigHash))

	for range ticker.C {
		if err := c.BaseConfig.Viper.ReadRemoteConfig(); err != nil {
			log.Error("Failed to read remote configuration", log.Err(err))
			continue
		}

		currentHash := hot_cfg.CalculateConfigHash(c.BaseConfig.Viper)

		// Compare hash values to detect changes
		if currentHash != lastConfigHash {
			log.Info("Configuration change detected", log.Str("lastConfigHash", lastConfigHash), log.Str("currentHash", currentHash))
			lastConfigHash = currentHash

			if err := c.UnmarshalConfig(); err != nil {
				log.Error("Failed to parse updated configuration", log.Err(err))
				continue
			}

			log.Info("Configuration manually updated", log.Any("config", NetCfg))
		}
	}
}
