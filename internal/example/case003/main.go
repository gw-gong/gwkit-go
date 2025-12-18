package main

import (
	"context"
	"time"

	"github.com/gw-gong/gwkit-go/hotcfg"
	"github.com/gw-gong/gwkit-go/internal/example/case003/config"

	// "github.com/gw-gong/gwkit-go/internal/examples/case003/netconfig"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/util"
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	util.ExitOnErr(context.Background(), err)
	defer syncFn()

	hlm := hotcfg.NewHotLoaderManager()

	localConfigOption := &hotcfg.LocalConfigOption{
		FilePath: "config",
		FileName: "config-dev.yaml",
		FileType: "yaml",
	}
	localConfig, err := config.NewConfig(localConfigOption)
	util.ExitOnErr(context.Background(), err)
	err = hlm.RegisterHotLoader(localConfig)
	util.ExitOnErr(context.Background(), err)

	// consulConfigOption := &hotcfg.ConsulConfigOption{
	// 	ConsulAddr: "127.0.0.1:8500",
	// 	ConsulKey:  "config/config-dev.yaml",
	// 	ConfigType: "yaml",
	// 	ReloadTime: 10,
	// }
	// consulConfig, err := netconfig.NewConfig(consulConfigOption)
	// util.ExitOnErr(context.Background(), err)
	// err = hlm.RegisterHotLoader(consulConfig)
	// util.ExitOnErr(context.Background(), err)

	util.ExitOnErr(context.Background(), hlm.Watch())

	// test
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

// func testNetConfig(config *netconfig.Config) {
// 	go func() {
// 		for {
// 			log.Info("testNetConfig", log.Any("config", config))
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()
// }
