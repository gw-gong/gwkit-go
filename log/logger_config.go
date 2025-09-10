package log

const (
	DefaultOutputFilePath string = "logs/app.log"

	OutputEncodingJSON    string = "json"
	OutputEncodingConsole string = "console"
)

type LoggerConfig struct {
	Level           string                `yaml:"level" json:"level"`                         // Log level
	OutputToFile    OutputToFileConfig    `yaml:"output_to_file" json:"output_to_file"`       // Output to file configuration
	OutputToConsole OutputToConsoleConfig `yaml:"output_to_console" json:"output_to_console"` // Output to console configuration
	AddCaller       bool                  `yaml:"add_caller" json:"add_caller"`               // Whether to add caller information
	StackTrace      StackTraceConfig      `yaml:"stack_trace" json:"stack_trace"`             // Stack trace configuration
}

type OutputToFileConfig struct {
	Enable     bool   `yaml:"enable" json:"enable"`           // Whether to enable output
	FilePath   string `yaml:"path" json:"path"`               // File path
	WithBuffer bool   `yaml:"with_buffer" json:"with_buffer"` // Whether to use buffer, will not be immediately flushed to the file
	MaxSize    int    `yaml:"max_size" json:"max_size"`       // Maximum file size (MB)
	MaxBackups int    `yaml:"max_backups" json:"max_backups"` // Maximum number of backup files
	MaxAge     int    `yaml:"max_age" json:"max_age"`         // Maximum retention days
	Compress   bool   `yaml:"compress" json:"compress"`       // Whether to compress after rotation
}

type OutputToConsoleConfig struct {
	Enable   bool   `yaml:"enable" json:"enable"`     // Whether to enable output
	Encoding string `yaml:"encoding" json:"encoding"` // Encoding
}

type StackTraceConfig struct {
	Enable     bool   `yaml:"enable" json:"enable"` // Whether to enable stack trace
	TraceLevel string `yaml:"level" json:"level"`   // Stack trace level
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
		StackTrace: StackTraceConfig{
			Enable:     false,
			TraceLevel: LoggerLevelError,
		},
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

	// Merge stack trace config
	if config.StackTrace.Enable {
		mergedConfig.StackTrace = config.StackTrace
	} else {
		// Merge level settings even if not enabled
		if config.StackTrace.TraceLevel != "" {
			mergedConfig.StackTrace.TraceLevel = config.StackTrace.TraceLevel
		}
	}

	return mergedConfig
}
