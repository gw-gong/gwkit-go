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
