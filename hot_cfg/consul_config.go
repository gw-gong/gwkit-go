package hot_cfg

type ConsulConfig interface {
	GetConsulReloadTime() int
	ReadConsulConfig() error
	CalculateConsulConfigHash() string
}

type consulConfig struct {
	enable     bool
	ConsulAddr string
	ConsulKey  string
	ConfigType string
	ReloadTime int // second
}