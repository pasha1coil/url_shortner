package initialize

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func Redis(ctx context.Context, config *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	status := rdb.Ping(ctx)
	if err := status.Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
