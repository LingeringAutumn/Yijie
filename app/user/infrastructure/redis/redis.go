package redis

import (
	"github.com/LingeringAutumn/Yijie/app/user/domain/repository"
	"github.com/redis/go-redis/v9"
)

type userRedis struct {
	client *redis.Client
}

func NewUserRedis(client *redis.Client) repository.UserRedis {
	cli := userRedis{client: client}
	return &cli
}
