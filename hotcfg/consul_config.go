package hotcfg

type ConsulConfig interface {
	GetConsulReloadTime() int
	ReadConsulConfig() error
	CalculateConsulConfigHash() string
}

type ConsulConfigOption struct {
	ConsulAddr string `json:"consulAddr" yaml:"consulAddr" mapstructure:"consulAddr"`
	ConsulKey  string `json:"consulKey" yaml:"consulKey" mapstructure:"consulKey"`
	ConfigType string `json:"configType" yaml:"configType" mapstructure:"configType"`
	ReloadTime int    `json:"reloadTime" yaml:"reloadTime" mapstructure:"reloadTime"` // second
}
