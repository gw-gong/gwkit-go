package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type corsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           time.Duration
}

type corsOptions func(config *corsConfig)

func WithCorsOptAllowOrigins(allowOrigins []string) corsOptions {
	return func(config *corsConfig) {
		config.AllowOrigins = allowOrigins
	}
}

func WithCorsOptAllowMethods(allowMethods []string) corsOptions {
	return func(config *corsConfig) {
		config.AllowMethods = allowMethods
	}
}

func WithCorsOptAllowHeaders(allowHeaders []string) corsOptions {
	return func(config *corsConfig) {
		config.AllowHeaders = allowHeaders
	}
}

func WithCorsOptAllowCredentials(allowCredentials bool) corsOptions {
	return func(config *corsConfig) {
		config.AllowCredentials = allowCredentials
	}
}

func WithCorsOptMaxAge(maxAge time.Duration) corsOptions {
	return func(config *corsConfig) {
		config.MaxAge = maxAge
	}
}

func Cors(options ...corsOptions) gin.HandlerFunc {
	config := corsConfig{
		AllowOrigins:     []string{},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	for _, option := range options {
		option(&config)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.AllowOrigins
	corsConfig.AllowMethods = config.AllowMethods
	corsConfig.AllowHeaders = config.AllowHeaders
	corsConfig.AllowCredentials = config.AllowCredentials
	corsConfig.MaxAge = config.MaxAge
	return cors.New(corsConfig)
}
