package service

import (
	redis_ "billiards/pkg/redis"
	"github.com/go-redis/redis"
)

type redisService struct {
	redis *redis.Client
}

var RedisService = &redisService{
	redis: redis_.GetRedis(),
}

func (r *redisService) Query(method, query string) (val interface{}, err error) {
	switch method {
	case "GET":
		if val, err = r.redis.Get(query).Result(); err != nil {
			return
		}
	case "LLEN":
		if val, err = r.redis.LLen(query).Result(); err != nil {
			return
		}
	case "TTL":
		if val, err = r.redis.TTL(query).Result(); err != nil {
			return
		}
	}

	return
}
