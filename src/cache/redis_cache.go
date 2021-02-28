package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	apierror "github.com/michelaquino/golang_api_skeleton/src/api_errors"
	"github.com/michelaquino/golang_api_skeleton/src/log"

	goredis "github.com/go-redis/redis/v8"
)

var (
	logger = log.GetLogger()
)

// Config has the redis configs
type Config struct {
	Name                  string
	Topology              string
	Host                  string
	Port                  int
	Password              string
	SentinelMasterName    string
	PoolSize              int
	MinIdleConnections    int
	PoolTimeout           time.Duration
	ConnectionIdleTimeout time.Duration
	DialTimeout           time.Duration
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	MaxRetries            int
	MaxRetryBackoff       time.Duration
}

// RedisCache is a redis cache object.
type RedisCache struct {
	client goredis.UniversalClient
}

// NewRedis instantiate a new redis cluster client
func NewRedis(config Config) *RedisCache {
	addresses := config.getAddresses()
	var universalClient goredis.UniversalClient

	switch config.Topology {
	case "cluster":
		universalClient = goredis.NewClusterClient(&goredis.ClusterOptions{
			Addrs:    addresses,
			Password: config.Password,

			PoolSize:    config.PoolSize,
			PoolTimeout: config.PoolTimeout,

			MinIdleConns: config.MinIdleConnections,
			IdleTimeout:  config.ConnectionIdleTimeout,

			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,

			MaxRetries:      config.MaxRetries,
			MaxRetryBackoff: config.MaxRetryBackoff,
		})
	case "sentinel":
		universalClient = goredis.NewFailoverClient(&goredis.FailoverOptions{
			MasterName: config.SentinelMasterName,

			SentinelAddrs: addresses,
			Password:      config.Password,

			PoolSize:    config.PoolSize,
			PoolTimeout: config.PoolTimeout,

			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,

			MaxRetries:      config.MaxRetries,
			MaxRetryBackoff: config.MaxRetryBackoff,

			MinIdleConns: config.MinIdleConnections,
			IdleTimeout:  config.ConnectionIdleTimeout,
		})
	default:
		var redisAddress string
		if len(addresses) > 0 {
			redisAddress = addresses[0]
		}

		universalClient = goredis.NewClient(&goredis.Options{
			Addr:     redisAddress,
			Password: config.Password,

			PoolSize:    config.PoolSize,
			PoolTimeout: config.PoolTimeout,

			MinIdleConns: config.MinIdleConnections,
			IdleTimeout:  config.ConnectionIdleTimeout,

			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		})
	}

	return &RedisCache{client: universalClient}
}

// Get is a method that gets a value from cache.
func (r RedisCache) Get(ctx context.Context, key string) (string, error) {
	cacheValue, err := r.client.Get(ctx, key).Result()

	logAction := fmt.Sprintf("get key %s", key)
	if err == goredis.Nil {
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
	err := r.client.Set(ctx, key, value, expire).Err()
	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	logger.Debug(ctx, logAction, "success", nil)
	return nil
}

func (c Config) getAddresses() []string {
	addresses := strings.Split(c.Host, ",")
	for i, address := range addresses {
		addresses[i] = fmt.Sprintf("%s:%d", address, c.Port)
	}

	return addresses
}
