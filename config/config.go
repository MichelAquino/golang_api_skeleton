package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	logEncodingKeyValue string = "key-value"
	logEncodingJSON     string = "json"
)

// Init application config
func Init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("api.host.port", 8888)
	viper.SetDefault("api.handler.timeout", "500ms")

	viper.SetDefault("log.level", "ERROR")
	viper.SetDefault("log.file.name", "")
	viper.SetDefault("log.to.file", false)
	viper.SetDefault("log.encoding", logEncodingJSON)

	viper.SetDefault("new_relic.is.enabled", false)
	viper.SetDefault("new_relic.licence.key", "")
	viper.SetDefault("new_relic.proxy.url", "")

	viper.SetDefault("redis.url", "redis")
	viper.SetDefault("redis.password", "")

	viper.SetDefault("mongo.uri", "mongodb://mongodb:27017")
}
