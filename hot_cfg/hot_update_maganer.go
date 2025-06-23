package hot_cfg

import (
	"fmt"
	"sync"
)

type HotUpdate interface {
	UnmarshalConfig() error
	Watch()
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
	if err := hotUpdate.UnmarshalConfig(); err != nil {
		return fmt.Errorf("hot update init failed: %w", err)
	}
	m.hotUpdates = append(m.hotUpdates, hotUpdate)
	return nil
}

// After all the registrations are completed, start the watch.
func (m *hotUpdateConfigManager) WatchAll() {
	for _, hotUpdate := range m.hotUpdates {
		hotUpdate.Watch()
	}
}
