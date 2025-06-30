package hot_cfg

import (
	"fmt"
	"sync"
	"time"

	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"
)

type HotUpdate interface {
	GetBaseConfig() *BaseConfig
	AsLocalConfig() LocalConfig
	AsConsulConfig() ConsulConfig

	// Only need to implement these two methods, others inherit from BaseConfig
	LoadConfig()
}

var (
	hucm *hotUpdateConfigManager
	once sync.Once
)

type hotUpdateConfigManager struct {
	hotUpdates []HotUpdate
}

func initHotUpdateManager() {
	once.Do(func() {
		hucm = &hotUpdateConfigManager{
			hotUpdates: make([]HotUpdate, 0),
		}
	})
}

func GetHotUpdateManager() *hotUpdateConfigManager {
	if hucm == nil {
		initHotUpdateManager()
	}
	return hucm
}

func (m *hotUpdateConfigManager) RegisterHotUpdateConfig(hotUpdate HotUpdate) error {
	hotUpdate.LoadConfig()
	m.hotUpdates = append(m.hotUpdates, hotUpdate)
	return nil
}

// After all the registrations are completed, start the watch.
func (m *hotUpdateConfigManager) Watch() error {
	var errors []error

	for _, hotUpdate := range m.hotUpdates {
		if localConfig := hotUpdate.AsLocalConfig(); localConfig != nil {
			localConfig.WatchLocalConfig(hotUpdate.LoadConfig)
		} else if consulConfig := hotUpdate.AsConsulConfig(); consulConfig != nil {
			go gwkit_common.WithRecover(func() {
				func(consulConfig ConsulConfig, hotUpdate HotUpdate) {
					m.watchConsulConfig(consulConfig, hotUpdate.LoadConfig)
				}(consulConfig, hotUpdate)
			})
		} else {
			errors = append(errors, fmt.Errorf("hot update config struct error: %v", hotUpdate.GetBaseConfig()))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("hot update config struct error: %v", errors)
	}
	return nil
}

func (m *hotUpdateConfigManager) watchConsulConfig(consulConfig ConsulConfig, loadConfig func()) {
	ticker := time.NewTicker(time.Duration(consulConfig.GetConsulReloadTime()) * time.Second)
	defer ticker.Stop()

	lastConfigHash := consulConfig.CalculateConsulConfigHash()

	for range ticker.C {
		if err := consulConfig.ReadConsulConfig(); err != nil {
			continue
		}

		currentConfigHash := consulConfig.CalculateConsulConfigHash()
		if currentConfigHash != lastConfigHash {
			lastConfigHash = currentConfigHash
			loadConfig()
		}
	}
}
