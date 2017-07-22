package context

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// LogConfig represents the log configuration
type LogConfig struct {
	LogLevel    string
	LogToFile   bool
	LogFileName string
}

// MongoConfig represents the MongoDB configuration
type MongoConfig struct {
	Address      string
	DatabaseName string
	Timeout      time.Duration
	Username     string
	Password     string
}

// RedisConfig represents the Redis configuration
type RedisConfig struct {
	RedisURL      string
	RedisPassword string
}

// MySQLConfig represents the MySQL configuration
type MySQLConfig struct {
	URL      string
	User     string
	Password string
	Database string
}

// APIConfig represents the API configuration
type APIConfig struct {
	LogConfig     *LogConfig
	MongoDBConfig *MongoConfig
	RedisConfig   *RedisConfig
	MySQLConfig   *MySQLConfig
}

var apiConfig *APIConfig
var onceConfig sync.Once

// GetAPIConfig return the instance of the APIConfig
func GetAPIConfig() *APIConfig {
	onceConfig.Do(func() {
		apiConfig = &APIConfig{
			LogConfig:     getLogConfig(),
			MongoDBConfig: getMongoConfig(),
			RedisConfig:   getRedisConfig(),
			MySQLConfig:   getMySQLConfig(),
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

func getMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		URL:      os.Getenv("MYSQL_URL"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}
}
