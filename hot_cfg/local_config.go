package hot_cfg

type LocalConfig interface {
	WatchLocalConfig(loadConfig func())
}

type localConfig struct {
	enable   bool
	FilePath string
	FileName string
	FileType string
}
