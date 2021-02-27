package cache

import (
	"context"
	"fmt"
	"time"

	apierror "github.com/michelaquino/golang_api_skeleton/src/api_errors"
	"github.com/michelaquino/golang_api_skeleton/src/log"
	"github.com/spf13/viper"

	"github.com/go-redis/redis"
)

var (
	logger = log.GetLogger()
)

// RedisCache is a redis cache object.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache returns a new instance of RedisCache.
func NewRedisCache() *RedisCache {
	var redisClient *redis.Client
	ctx := context.Background()
	redisClient = redis.NewClient(&redis.Options{
		ReadTimeout:  time.Duration(1) * time.Second,
		WriteTimeout: time.Duration(1) * time.Second,
		Addr:         viper.GetString("redis.url"),
		Password:     viper.GetString("redis.password"),
		DB:           0,
		PoolSize:     5000,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		logger.Error(ctx, "ping Redis", err.Error(), nil)
	}

	return &RedisCache{
		client: redisClient,
	}
}

// Get is a method that gets a value from cache.
func (r RedisCache) Get(ctx context.Context, key string) (string, error) {
	cacheValue, err := r.client.Get(key).Result()

	logAction := fmt.Sprintf("get key %s", key)
	if err == redis.Nil {
		logger.Info(ctx, logAction, "", nil)
		return "", apierror.ErrNotFoundOnCache
	}

	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return "", apierror.ErrGetCacheValue
	}

	logger.Info(ctx, logAction, "success", nil)
	return cacheValue, nil
}

// Set is a method that sets a value to cache.
func (r RedisCache) Set(ctx context.Context, key, value string, expireInSec int) error {
	expire := time.Duration(expireInSec) * time.Second

	logAction := fmt.Sprintf("get key %s with expiration %d", key, expireInSec)
	err := r.client.Set(key, value, expire).Err()
	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	logger.Debug(ctx, logAction, "success", nil)
	return nil
}
