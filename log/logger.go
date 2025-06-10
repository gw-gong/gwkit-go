package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GlobalLogger() *zap.Logger {
	return zap.L()
}

func InitGlobalLogger(loggerConfig *LoggerConfig) error {
	logger, err := newLogger(loggerConfig)
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}

func newLogger(loggerConfig *LoggerConfig) (*zap.Logger, error) {
	loggerConfig = MergeCfgIntoDefault(loggerConfig)

	zapCfg := zap.NewProductionConfig()

	zapCfg.Level.SetLevel(MapLoggerLevel(loggerConfig.Level))

	// Create basic encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create encoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create WriteSyncer
	var writeSyncer zapcore.WriteSyncer
	if loggerConfig.OutputToFile.Enable && loggerConfig.OutputToFile.FilePath != "" {
		// Ensure directory exists
		dir := filepath.Dir(loggerConfig.OutputToFile.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Use lumberjack to split logs
		lumber := &lumberjack.Logger{
			Filename:   loggerConfig.OutputToFile.FilePath,
			MaxSize:    loggerConfig.OutputToFile.MaxSize,    // Unit: MB
			MaxBackups: loggerConfig.OutputToFile.MaxBackups, // Maximum number of old files to retain
			MaxAge:     loggerConfig.OutputToFile.MaxAge,     // Maximum number of days to retain old files
			Compress:   loggerConfig.OutputToFile.Compress,   // Whether to compress after rotation
		}
		writeSyncer = zapcore.AddSync(lumber)
	} else if loggerConfig.OutputToConsole.Enable && IsSupportedEncodingType(loggerConfig.OutputToConsole.Encoding) {
		writeSyncer = zapcore.AddSync(os.Stdout)
		if loggerConfig.OutputToConsole.Encoding == OutputEncodingConsole {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}
	} else {
		return nil, errors.New("no valid output configured: either output_to_file or output_to_console must be enabled with valid settings")
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, zapCfg.Level)

	// Create Logger instance
	var logger *zap.Logger
	if loggerConfig.AddCaller {
		logger = zap.New(core, zap.AddCaller())
	} else {
		logger = zap.New(core)
	}

	// Add stack trace
	if loggerConfig.StackTrace.Enable {
		stackLevel := MapLoggerLevel(loggerConfig.StackTrace.TraceLevel)
		logger = logger.WithOptions(zap.AddStacktrace(stackLevel))
	}

	return logger, nil
}
