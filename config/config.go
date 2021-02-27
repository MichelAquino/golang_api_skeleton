package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	continueWatching = "continue_watching"
)

// Init application config
func Init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("api.host.port", 8888)

	viper.SetDefault("log.level", "ERROR")
	viper.SetDefault("log.file.name", "")
	viper.SetDefault("log.to.file", false)

	viper.SetDefault("new_relic.is.enabled", false)
	viper.SetDefault("new_relic.licence.key", "")
	viper.SetDefault("new_relic.proxy.url", "")

	viper.SetDefault("redis.url", "redis")
	viper.SetDefault("redis.password", "")

	viper.SetDefault("mongo.url", "mongodb")
	viper.SetDefault("mongo.port", 27017)
	viper.SetDefault("mongo.username", "")
	viper.SetDefault("mongo.password", "")
	viper.SetDefault("mongo.database.name", "api")
	viper.SetDefault("mongo.timeout", "60s")
}
