package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type corsConfig struct {
	AllowOrigins     []string `json:"allow_origins" yaml:"allow_origins" mapstructure:"allow_origins"`
	AllowMethods     []string `json:"allow_methods" yaml:"allow_methods" mapstructure:"allow_methods"`
	AllowHeaders     []string `json:"allow_headers" yaml:"allow_headers" mapstructure:"allow_headers"`
	AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials" mapstructure:"allow_credentials"`
	MaxAge           int      `json:"max_age" yaml:"max_age" mapstructure:"max_age"`
}

func MergeDefaultCorsConfig(config *corsConfig) *corsConfig {
	// default config
	defaultConfig := corsConfig{
		AllowOrigins:     []string{},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * 60 * 60, // 12 hours
	}
	if config == nil {
		return &defaultConfig
	}
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = defaultConfig.AllowOrigins
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = defaultConfig.AllowMethods
	}
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = defaultConfig.AllowHeaders
	}
	if config.MaxAge == 0 {
		config.MaxAge = defaultConfig.MaxAge
	}
	return config
}

func Cors(config *corsConfig) gin.HandlerFunc {
	config = MergeDefaultCorsConfig(config)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.AllowOrigins
	corsConfig.AllowMethods = config.AllowMethods
	corsConfig.AllowHeaders = config.AllowHeaders
	corsConfig.AllowCredentials = config.AllowCredentials
	corsConfig.MaxAge = time.Duration(config.MaxAge) * time.Hour
	return cors.New(corsConfig)
}
