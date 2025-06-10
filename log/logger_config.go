package log

const (
	DefaultOutputFilePath string = "logs/app.log"

	OutputEncodingJSON    string = "json"
	OutputEncodingConsole string = "console"
)

type OutputToFileConfig struct {
	Enable     bool   `yaml:"enable" json:"enable"`           // Whether to enable output
	FilePath   string `yaml:"path" json:"path"`               // File path
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

type LoggerConfig struct {
	Level           string                `yaml:"level" json:"level"`                         // Log level
	OutputToFile    OutputToFileConfig    `yaml:"output_to_file" json:"output_to_file"`       // Output to file configuration
	OutputToConsole OutputToConsoleConfig `yaml:"output_to_console" json:"output_to_console"` // Output to console configuration
	AddCaller       bool                  `yaml:"add_caller" json:"add_caller"`               // Whether to add caller information
	StackTrace      StackTraceConfig      `yaml:"stack_trace" json:"stack_trace"`             // Stack trace configuration
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

func IsSupportedEncodingType(encoding string) bool {
	return encoding == OutputEncodingJSON || encoding == OutputEncodingConsole
}
