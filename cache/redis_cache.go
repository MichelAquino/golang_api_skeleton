package cache

import (
	"strconv"
	"time"

	"github.com/michelaquino/golang_api_skeleton/api_errors"
	"github.com/michelaquino/golang_api_skeleton/context"

	"github.com/go-redis/redis"
)

// RedisCache is a redis cache object.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache returns a new instance of RedisCache.
func NewRedisCache() *RedisCache {
	cacheLogger := context.GetLogger()
	apiConfig := context.GetAPIConfig()

	var redisClient *redis.Client
	cacheLogger.Info("NewRedisCache", "Constructor", "", "", "", "", "")
	redisClient = redis.NewClient(&redis.Options{
		ReadTimeout:  time.Duration(1) * time.Second,
		WriteTimeout: time.Duration(1) * time.Second,
		Addr:         apiConfig.RedisConfig.RedisURL,
		Password:     apiConfig.RedisConfig.RedisPassword,
		DB:           0,
		PoolSize:     5000,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		cacheLogger.Error("NewRedisCache", "Constructor", "", "", "Ping Redis", "Error", err.Error())
	}

	return &RedisCache{
		client: redisClient,
	}
}

// Get is a method that gets a value from cache.
func (r RedisCache) Get(key string) (string, error) {
	cacheLogger := context.GetLogger()

	cacheValue, err := r.client.Get(key).Result()
	if err == redis.Nil {
		cacheLogger.Debug("RedisCache", "Get", "", "", "Get key", "Key "+key+" not found on cache", redis.Nil.Error())
		return "", apierror.ErrNotFoundOnCache
	}

	if err != nil {
		cacheLogger.Error("RedisCache", "Get", "", "", "Get key", "Error", err.Error())
		return "", apierror.ErrGetCacheValue
	}

	cacheLogger.Error("RedisCache", "Get", "", "", "Get key", "Success", "Object getted with success")
	return cacheValue, nil
}

// Set is a method that sets a value to cache.
func (r RedisCache) Set(key, value string, expireInSec int) error {
	cacheLogger := context.GetLogger()

	expire := time.Duration(expireInSec) * time.Second

	err := r.client.Set(key, value, expire).Err()
	if err != nil {
		cacheLogger.Error("RedisCache", "Set", "", "", "Set key/value on cache", "Error", err.Error())
	}

	cacheLogger.Debug("RedisCache", "Set", "", "", "Set key/value on cache", "Success", "Set key "+key+" on cache with expire in "+strconv.Itoa(expireInSec)+" seconds")
	return err
}
