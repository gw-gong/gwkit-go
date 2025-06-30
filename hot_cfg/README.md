# Hot Configuration

这个目录提供了热配置功能，支持本地文件和Consul两种配置源的动态更新，无需重启应用即可生效配置变更。

## 功能特性

- **多配置源支持**: 支持本地文件和Consul远程配置
- **热更新**: 配置变更时自动重新加载，无需重启应用
- **线程安全**: 使用读写锁保证配置访问的线程安全
- **配置哈希检测**: 通过MD5哈希检测配置变更，避免不必要的重载
- **统一管理**: 提供统一的热更新管理器，支持多个配置同时监控

## 文件说明

### base_config.go

提供基础配置功能：

- `BaseConfig`: 基础配置结构体，包含Viper实例和配置选项
- `NewBaseConfig()`: 创建基础配置实例，支持本地文件和Consul配置
- `WithLocalConfig()`: 本地文件配置选项
- `WithConsulConfig()`: Consul远程配置选项
- `WatchLocalConfig()`: 监控本地文件变更
- `ReadConsulConfig()`: 读取Consul远程配置
- `CalculateConsulConfigHash()`: 计算配置哈希值

### hot_update_manager.go

提供热更新管理器：

- `HotUpdate`: 热更新接口，定义配置更新方法
- `hotUpdateConfigManager`: 热更新管理器，管理多个配置实例
- `GetHotUpdateManager()`: 获取单例热更新管理器
- `RegisterHotUpdateConfig()`: 注册热更新配置
- `Watch()`: 启动所有配置的监控

### local_config.go

本地配置接口定义：

- `LocalConfig`: 本地配置接口
- `WatchLocalConfig()`: 监控本地配置文件变更

### consul_config.go

Consul配置接口定义：

- `ConsulConfig`: Consul配置接口
- `GetConsulReloadTime()`: 获取重载时间间隔
- `ReadConsulConfig()`: 读取Consul配置
- `CalculateConsulConfigHash()`: 计算配置哈希

### utils.go

工具函数：

- `CalculateConfigHash()`: 计算配置的MD5哈希值，用于检测配置变更

## 使用示例

### 本地文件配置

```go
package config

import (
	"fmt"
	"sync"
	"time"

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

func (c *Config) LoadConfig() {
	c.BaseConfig.Mu.Lock()
	defer c.BaseConfig.Mu.Unlock()

	if err := c.BaseConfig.Viper.Unmarshal(&c); err != nil {
		log.Error("unmarshal config failed", log.Err(err))
		return
	}

	log.Info("LoadConfig", log.Any("config", c))
}

```

### Consul远程配置

```go
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
```

### 使用热更新管理器

```go
package main

import (
	"context"
	"time"

	"github.com/gw-gong/gwkit-go/hot_cfg"
	"github.com/gw-gong/gwkit-go/internal/examples/case003/config"
	// "github.com/gw-gong/gwkit-go/internal/examples/case003/net_config"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	gwkit_common.ExitOnErr(context.Background(), err)
	defer syncFn()

	hucm := hot_cfg.GetHotUpdateManager()

	err = config.InitConfig("config", "config-dev.yaml", "yaml")
	gwkit_common.ExitOnErr(context.Background(), err)
	err = hucm.RegisterHotUpdateConfig(config.Cfg)
	gwkit_common.ExitOnErr(context.Background(), err)

	// err = net_config.InitNetConfig("127.0.0.1:8500", "config/config-dev.yaml", "yaml", 10)
	// gwkit_common.ExitOnErr(context.Background(), err)
	// err = hucm.RegisterHotUpdateConfig(net_config.NetCfg)
	// gwkit_common.ExitOnErr(context.Background(), err)

	hucm.Watch()

	// 测试热更新
	testLocoalConfig()
	// testNetConfig()
	select {}
}

func testLocoalConfig() {
	go func() {
		for {
			log.Info("testLocoalConfig", log.Any("config", config.Cfg))
			time.Sleep(5 * time.Second)
		}
	}()
}

// func testNetConfig() {
// 	go func() {
// 		for {
// 			log.Info("testNetConfig", log.Any("config", net_config.NetCfg))
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()
// }

```

## 配置格式

### 本地配置文件 (config-dev.yaml)

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: password
api:
  key: api-key-123
  timeout: 30s
```

### Consul配置

在Consul中存储相同格式的配置内容，通过指定的key进行访问。

## 工作流程

1. **初始化配置**:

   - 创建配置结构体，嵌入 `BaseConfig`, 注意一定要使用label "mapstructure"
   - 使用 `NewBaseConfig()`初始化基础配置
   - 实现 `LoadConfig()`方法
2. **注册配置**:

   - 获取热更新管理器实例
   - 调用 `RegisterHotUpdateConfig()`注册配置
3. **启动监控**:

   - 调用 `Watch()`启动所有配置的监控
   - 本地文件：使用fsnotify监控文件变更
   - Consul配置：定时检查配置哈希值变化
4. **配置更新**:

   - 检测到配置变更时自动调用 `LoadConfig()`
   - 重新解析配置并更新内存中的配置对象
   - 记录更新日志

## 配置说明

- **本地文件监控**: 基于fsnotify实现文件变更监控
- **Consul监控**: 基于定时器检查配置哈希值变化
- **线程安全**: 使用读写锁保护配置访问
- **错误处理**: 配置更新失败时记录错误日志，不影响应用运行

## 依赖

- `github.com/spf13/viper`: 配置管理库
- `github.com/fsnotify/fsnotify`: 文件系统监控
- `github.com/spf13/viper/remote`: 远程配置支持
- `crypto/md5`: 配置哈希计算
