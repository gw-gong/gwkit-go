package hot_cfg

type LocalConfig interface {
	WatchLocalConfig(reloadConfig func())
}

type localConfig struct {
	enable   bool
	FilePath string
	FileName string
	FileType string
}
