package hot_cfg

type HotLoader interface {
	GetBaseConfig() *BaseConfig
	AsLocalConfig() LocalConfig
	AsConsulConfig() ConsulConfig

	// Only need to implement these two methods, others inherit from BaseConfig
	LoadConfig()
}
