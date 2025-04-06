package redis

import (
	"github.com/redis/go-redis/v9"

	"github.com/LingeringAutumn/Yijie/app/user/domain/repository"
)

type userRedis struct {
	client *redis.Client
}

func NewUserRedis(client *redis.Client) repository.UserRedis {
	cli := userRedis{client: client}
	return &cli
}
