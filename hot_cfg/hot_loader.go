package hot_cfg

type HotLoader interface {
	BaseConfigCapable
	LoadConfig()
}
