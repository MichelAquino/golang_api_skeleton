package context

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// LogConfig represents Log configuration.
type LogConfig struct {
	LogLevel    string
	LogToFile   bool
	LogFileName string
}

// MongoConfig represents MongoDB configuration.
type MongoConfig struct {
	Address      string
	DatabaseName string
	Timeout      time.Duration
	Username     string
	Password     string
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	RedisURL      string
	RedisPassword string
}

// APIConfig represents API configuration.
type APIConfig struct {
	LogConfig     *LogConfig
	MongoDBConfig *MongoConfig
	RedisConfig   *RedisConfig
	HostPort      int
}

var apiConfig *APIConfig
var onceConfig sync.Once

// GetAPIConfig returns the instance of an APIConfig.
func GetAPIConfig() *APIConfig {
	onceConfig.Do(func() {
		apiConfig = &APIConfig{
			LogConfig:     getLogConfig(),
			MongoDBConfig: getMongoConfig(),
			RedisConfig:   getRedisConfig(),
			HostPort:      getAPIHostPort(),
		}
	})

	return apiConfig
}

func getLogConfig() *LogConfig {
	return &LogConfig{
		LogLevel:    os.Getenv("LOG_LEVEL"),
		LogToFile:   getLogToFileFlag(),
		LogFileName: os.Getenv("LOG_FILE_NAME"),
	}
}

func getLogToFileFlag() bool {
	if logToFile, err := strconv.ParseBool(os.Getenv("LOG_TO_FILE")); err == nil {
		return logToFile
	}

	return false
}

func getRedisConfig() *RedisConfig {
	return &RedisConfig{
		RedisURL:      os.Getenv("REDIS_URL"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}

func getMongoConfig() *MongoConfig {
	mongoURL := os.Getenv("MONGO_URL")
	mongoPort := getMongoPort()

	mongoAddress := fmt.Sprintf("%s:%d", mongoURL, mongoPort)
	mongoDatabaseName := os.Getenv("MONGO_DATABASE_NAME")
	mongoTimeout := getMongoTimeout()
	mongoUserName := os.Getenv("MONGO_DATABASE_USERNAME")
	mongoPassword := os.Getenv("MONGO_DATABASE_PASSWORD")

	return &MongoConfig{
		Address:      mongoAddress,
		DatabaseName: mongoDatabaseName,
		Timeout:      mongoTimeout,
		Username:     mongoUserName,
		Password:     mongoPassword,
	}
}

func getMongoPort() int {
	mongoPort, err := strconv.Atoi(os.Getenv("MONGO_PORT"))
	if err != nil {
		return 27017
	}

	return mongoPort
}

func getMongoTimeout() time.Duration {
	mongoTimeout, err := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
	if err != nil {
		return time.Duration(60) * time.Second
	}

	return time.Duration(mongoTimeout) * time.Second
}

func getAPIHostPort() int {
	apiPort, err := strconv.Atoi(os.Getenv("HOST_PORT"))
	if err != nil {
		apiPort = 8888
	}

	return apiPort
}
