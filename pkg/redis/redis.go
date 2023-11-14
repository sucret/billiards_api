package redis

import (
	"gin-api/pkg/config"
	"sync"

	"github.com/go-redis/redis"
)

var (
	once     sync.Once
	instance *redis.Client
)

func NewRedis() *redis.Client {
	redisConfig := config.GetConfig().Redis

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + redisConfig.Port,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.DB,       // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		// global.Log.Error("Redis connect ping failed, err:", zap.Any("err", err))
		return nil
	}
	return client
}

func GetRedis() *redis.Client {
	once.Do(func() {
		instance = NewRedis()
	})

	return instance
}
