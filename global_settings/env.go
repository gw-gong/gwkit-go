package global_settings

const (
	ENV_DEV  = "dev"
	ENV_TEST = "test"

	ENV_STAGING = "staging"

	ENV_PROD = "prod"
	ENV_LIVE = "live"
)

var globalEnv = ENV_DEV

func SetEnv(env string) {
	globalEnv = env
}

func GetEnv() string {
	return globalEnv
}
