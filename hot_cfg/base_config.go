package hot_cfg

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type BaseConfig struct {
	Viper        *viper.Viper
	Mu           sync.RWMutex
	LocalConfig  localConfig
	ConsulConfig consulConfig
}

type option func(*BaseConfig)

func WithLocalConfig(filePath, fileName, fileType string) option {
	return func(c *BaseConfig) {
		c.LocalConfig.enable = true
		c.LocalConfig.FilePath = filePath
		c.LocalConfig.FileName = fileName
		c.LocalConfig.FileType = fileType
	}
}

func WithConsulConfig(consulAddr, consulKey, configType string, reloadTime int) option {
	return func(c *BaseConfig) {
		c.ConsulConfig.enable = true
		c.ConsulConfig.ConsulAddr = consulAddr
		c.ConsulConfig.ConsulKey = consulKey
		c.ConsulConfig.ConfigType = configType
		c.ConsulConfig.ReloadTime = reloadTime
	}
}

func NewBaseConfig(opts ...option) (*BaseConfig, error) {
	c := &BaseConfig{
		Viper: viper.New(),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.LocalConfig.enable && c.ConsulConfig.enable {
		return nil, fmt.Errorf("localConfig and consulConfig cannot be enabled at the same time")
	}

	if !c.LocalConfig.enable && !c.ConsulConfig.enable {
		return nil, fmt.Errorf("localConfig or consulConfig must be enabled")
	}

	if c.LocalConfig.enable && c.LocalConfig.FilePath != "" && c.LocalConfig.FileName != "" && c.LocalConfig.FileType != "" {
		c.Viper.SetConfigType(c.LocalConfig.FileType)
		c.Viper.SetConfigName(c.LocalConfig.FileName)
		c.Viper.AddConfigPath(c.LocalConfig.FilePath)
		if err := c.Viper.ReadInConfig(); err != nil {
			return nil, err
		}
	} else if c.ConsulConfig.enable && c.ConsulConfig.ConsulAddr != "" && c.ConsulConfig.ConsulKey != "" {
		if err := c.Viper.AddRemoteProvider("consul", c.ConsulConfig.ConsulAddr, c.ConsulConfig.ConsulKey); err != nil {
			return nil, err
		}
		c.Viper.SetConfigType(c.ConsulConfig.ConfigType)
		if err := c.Viper.ReadRemoteConfig(); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("localConfig or consulConfig must be enabled, and you need to prepare the relevant parameters")
	}

	return c, nil
}

func (c *BaseConfig) GetBaseConfig() *BaseConfig {
	return c
}

func (c *BaseConfig) AsLocalConfig() LocalConfig {
	if c.LocalConfig.enable {
		return c
	}
	return nil
}

func (c *BaseConfig) WatchLocalConfig(loadConfig func()) {
	if c.LocalConfig.enable {
		c.Viper.WatchConfig()
		c.Viper.OnConfigChange(func(e fsnotify.Event) {
			loadConfig()
		})
	}
}

func (c *BaseConfig) AsConsulConfig() ConsulConfig {
	if c.ConsulConfig.enable {
		return c
	}
	return nil
}

func (c *BaseConfig) GetConsulReloadTime() int {
	return c.ConsulConfig.ReloadTime
}

func (c *BaseConfig) ReadConsulConfig() error {
	if err := c.Viper.ReadRemoteConfig(); err != nil {
		return fmt.Errorf("failed to read remote configuration: %w, consulAddr: %s, consulKey: %s, configType: %s",
			err, c.ConsulConfig.ConsulAddr, c.ConsulConfig.ConsulKey, c.ConsulConfig.ConfigType)
	}
	return nil
}

func (c *BaseConfig) CalculateConsulConfigHash() string {
	return CalculateConfigHash(c.Viper)
}
