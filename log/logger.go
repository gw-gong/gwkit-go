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

func InitGlobalLogger(loggerConfig *LoggerConfig) (func(), error) {
	logger, syncGlobalLogger, err := newLogger(loggerConfig)
	if err != nil {
		return syncGlobalLogger, err
	}

	zap.ReplaceGlobals(logger)
	return syncGlobalLogger, nil
}

func newLogger(loggerConfig *LoggerConfig) (*zap.Logger, func(), error) {
	// Sync global logger
	var syncGlobalLogger = func() {
		_ = zap.L().Sync()
	}

	loggerConfig = mergeCfgIntoDefault(loggerConfig)

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
			return nil, syncGlobalLogger, fmt.Errorf("failed to create log directory: %w", err)
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

		syncGlobalLogger = func() {
			if err := zap.L().Sync(); err != nil {
				if f, err := os.OpenFile(loggerConfig.OutputToFile.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
					defer f.Close()
					fmt.Fprintf(f, "Failed to sync global logger: %v\n", err)
				}
			}
		}
	} else if loggerConfig.OutputToConsole.Enable && IsSupportedEncodingType(loggerConfig.OutputToConsole.Encoding) {
		writeSyncer = zapcore.AddSync(os.Stdout)
		if loggerConfig.OutputToConsole.Encoding == OutputEncodingConsole {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		syncGlobalLogger = func() {
			fmt.Println("Output to console, no sync needed, please ignore, no need to modify your code")
		}
	} else {
		return nil, syncGlobalLogger, errors.New("no valid output configured: either output_to_file or output_to_console must be enabled with valid settings")
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, zapCfg.Level)

	// Create Logger options
	var loggerOptions []zap.Option
	if loggerConfig.AddCaller {
		loggerOptions = append(loggerOptions, zap.AddCaller())
		// add caller skip
		loggerOptions = append(loggerOptions, zap.AddCallerSkip(1)) // 由于所有的log函数已封装一层，所以需要跳过一层
	}

	// create logger
	logger := zap.New(core, loggerOptions...)

	// Add stack trace
	if loggerConfig.StackTrace.Enable {
		stackLevel := MapLoggerLevel(loggerConfig.StackTrace.TraceLevel)
		logger = logger.WithOptions(zap.AddStacktrace(stackLevel))
	}

	return logger, syncGlobalLogger, nil
}
