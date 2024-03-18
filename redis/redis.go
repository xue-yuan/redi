package redis

import (
	"context"
	"redi/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Initialize() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisAddr,
		Password: config.Config.RedisPassword,
		DB:       config.Config.RedisDB,
	})

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}
