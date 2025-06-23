package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gw-gong/gwkit-go/hot_cfg"
	"github.com/gw-gong/gwkit-go/log"
)

var (
	Cfg  *Config
	once sync.Once
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

func InitConfig(filePath, fileName, fileType string) error {
	var err error
	once.Do(func() {
		Cfg = &Config{}
		Cfg.BaseConfig, err = hot_cfg.NewBaseConfig(
			hot_cfg.WithLocalConfig(filePath, fileName, fileType),
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
	c.BaseConfig.Viper.WatchConfig()
	c.BaseConfig.Viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed", log.Str("file", e.Name), log.Str("operation", e.Op.String()))

		// 重新解析配置
		if err := c.UnmarshalConfig(); err != nil {
			log.Error("Failed to reload config", log.Str("file", e.Name), log.Err(err))
			return
		}

		log.Info("Config reloaded successfully", log.Any("config", c))
	})
}
