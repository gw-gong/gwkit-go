package log

const (
	DefaultOutputFilePath string = "logs/app.log"

	OutputEncodingJSON    string = "json"
	OutputEncodingConsole string = "console"
)

type LoggerConfig struct {
	Level           string                `yaml:"level" json:"level" mapstructure:"level"`                                     // Log level
	OutputToFile    OutputToFileConfig    `yaml:"output_to_file" json:"output_to_file" mapstructure:"output_to_file"`          // Output to file configuration
	OutputToConsole OutputToConsoleConfig `yaml:"output_to_console" json:"output_to_console" mapstructure:"output_to_console"` // Output to console configuration
	AddCaller       bool                  `yaml:"add_caller" json:"add_caller" mapstructure:"add_caller"`                      // Whether to add caller information
}

type OutputToFileConfig struct {
	Enable     bool   `yaml:"enable" json:"enable" mapstructure:"enable"`                // Whether to enable output
	FilePath   string `yaml:"path" json:"path" mapstructure:"path"`                      // File path
	WithBuffer bool   `yaml:"with_buffer" json:"with_buffer" mapstructure:"with_buffer"` // Whether to use buffer, will not be immediately flushed to the file
	MaxSize    int    `yaml:"max_size" json:"max_size" mapstructure:"max_size"`          // Maximum file size (MB)
	MaxBackups int    `yaml:"max_backups" json:"max_backups" mapstructure:"max_backups"` // Maximum number of backup files
	MaxAge     int    `yaml:"max_age" json:"max_age" mapstructure:"max_age"`             // Maximum retention days
	Compress   bool   `yaml:"compress" json:"compress" mapstructure:"compress"`          // Whether to compress after rotation
}

type OutputToConsoleConfig struct {
	Enable   bool   `yaml:"enable" json:"enable" mapstructure:"enable"`       // Whether to enable output
	Encoding string `yaml:"encoding" json:"encoding" mapstructure:"encoding"` // Encoding
}

func IsSupportedEncodingType(encoding string) bool {
	return encoding == OutputEncodingJSON || encoding == OutputEncodingConsole
}

func NewDefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level: LoggerLevelDebug,
		OutputToFile: OutputToFileConfig{
			Enable:     false,
			FilePath:   DefaultOutputFilePath,
			WithBuffer: false,
			MaxSize:    500,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		},
		OutputToConsole: OutputToConsoleConfig{
			Enable:   true,
			Encoding: OutputEncodingConsole,
		},
		AddCaller: true,
	}
}

func mergeCfgIntoDefault(config *LoggerConfig) *LoggerConfig {
	if config == nil {
		return NewDefaultLoggerConfig()
	}

	// Create a new config object based on default config
	mergedConfig := NewDefaultLoggerConfig()

	// Merge log level
	if config.Level != "" {
		mergedConfig.Level = config.Level
	}

	// Merge file output config
	if config.OutputToFile.Enable {
		mergedConfig.OutputToFile = config.OutputToFile
	} else {
		// Merge other file config items even if not enabled
		if config.OutputToFile.FilePath != "" {
			mergedConfig.OutputToFile.FilePath = config.OutputToFile.FilePath
		}
		if config.OutputToFile.MaxSize > 0 {
			mergedConfig.OutputToFile.MaxSize = config.OutputToFile.MaxSize
		}
		if config.OutputToFile.MaxBackups > 0 {
			mergedConfig.OutputToFile.MaxBackups = config.OutputToFile.MaxBackups
		}
		if config.OutputToFile.MaxAge > 0 {
			mergedConfig.OutputToFile.MaxAge = config.OutputToFile.MaxAge
		}
		mergedConfig.OutputToFile.Compress = config.OutputToFile.Compress
	}

	// Merge console output config
	if config.OutputToConsole.Enable {
		mergedConfig.OutputToConsole = config.OutputToConsole
	} else {
		// Merge encoding settings even if not enabled
		if config.OutputToConsole.Encoding != "" {
			mergedConfig.OutputToConsole.Encoding = config.OutputToConsole.Encoding
		}
	}

	// Merge caller info config
	mergedConfig.AddCaller = config.AddCaller

	return mergedConfig
}
