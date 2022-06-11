package cache

import (
	"github.com/dhruv42/hraasaka/config"
	"github.com/go-redis/redis"
)

func ConnectCache() (*redis.Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.CacheUrl,
		Password: "",
		DB:       0,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
