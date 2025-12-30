package netconfig

import (
	"github.com/gw-gong/gwkit-go/hotcfg"
	"github.com/gw-gong/gwkit-go/log"
)

type Config struct {
	hotcfg.BaseConfigCapable
	Database struct {
		Host     string `yaml:"host" mapstructure:"host"`
		Port     int    `yaml:"port" mapstructure:"port"`
		Username string `yaml:"username" mapstructure:"username"`
		Password string `yaml:"password" mapstructure:"password"`
	} `yaml:"database" mapstructure:"database"`
	API struct {
		Key     string `yaml:"key" mapstructure:"key"`
		Timeout int    `yaml:"timeout" mapstructure:"timeout"`
	} `yaml:"api" mapstructure:"api"`
}

func (c *Config) LoadConfig() {
	if err := c.Unmarshal(&c); err != nil {
		log.Error("unmarshal config failed", log.Err(err))
		return
	}

	log.Info("LoadConfig", log.Any("config", c))
}

func NewConfig(consulConfigOption *hotcfg.ConsulConfigOption) (config *Config, err error) {
	config = &Config{}
	config.BaseConfigCapable, err = hotcfg.NewConsulBaseConfigCapable(consulConfigOption)
	config.LoadConfig()
	return config, err
}
