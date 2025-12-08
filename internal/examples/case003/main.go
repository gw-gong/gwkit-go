package main

import (
	"context"
	"time"

	"github.com/gw-gong/gwkit-go/hot_cfg"
	"github.com/gw-gong/gwkit-go/internal/examples/case003/config"

	// "github.com/gw-gong/gwkit-go/internal/examples/case003/net_config"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/util/common"
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	gwkit_common.ExitOnErr(context.Background(), err)
	defer syncFn()

	hlm := hot_cfg.NewHotLoaderManager()

	localConfigOption := &hot_cfg.LocalConfigOption{
		FilePath: "config",
		FileName: "config-dev.yaml",
		FileType: "yaml",
	}
	localConfig, err := config.NewConfig(localConfigOption)
	gwkit_common.ExitOnErr(context.Background(), err)
	err = hlm.RegisterHotLoader(localConfig)
	gwkit_common.ExitOnErr(context.Background(), err)

	// consulConfigOption := &hot_cfg.ConsulConfigOption{
	// 	ConsulAddr: "127.0.0.1:8500",
	// 	ConsulKey:  "config/config-dev.yaml",
	// 	ConfigType: "yaml",
	// 	ReloadTime: 10,
	// }
	// consulConfig, err := net_config.NewConfig(consulConfigOption)
	// gwkit_common.ExitOnErr(context.Background(), err)
	// err = hlm.RegisterHotLoader(consulConfig)
	// gwkit_common.ExitOnErr(context.Background(), err)

	gwkit_common.ExitOnErr(context.Background(), hlm.Watch())

	// 测试热更新
	testLocoalConfig(localConfig)
	// testNetConfig(consulConfig)
	select {}
}

func testLocoalConfig(config *config.Config) {
	go func() {
		for {
			log.Info("testLocoalConfig", log.Any("config", config))
			time.Sleep(5 * time.Second)
		}
	}()
}

// func testNetConfig(config *net_config.Config) {
// 	go func() {
// 		for {
// 			log.Info("testNetConfig", log.Any("config", config))
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()
// }
