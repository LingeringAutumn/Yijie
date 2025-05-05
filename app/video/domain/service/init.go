package service

import (
	"github.com/LingeringAutumn/Yijie/app/video/domain/repository"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

type VideoService struct {
	db    repository.VideoDB
	redis repository.VideoRedis
	rpc   repository.VideoRPC
	sf    *utils.Snowflake
}

func NewVideoService(db repository.VideoDB, redis repository.VideoRedis, sf *utils.Snowflake, rpc repository.VideoRPC) *VideoService {
	if db == nil {
		panic("videoService`s db should not be nil")
	}
	if redis == nil {
		panic("videoervice`s redis should not be nil")
	}
	if sf == nil {
		panic("videoService`s sf should not be nil")
	}
	if rpc == nil {
		panic("videoervice`s rpc should not be nil")
	}
	svc := &VideoService{
		db:    db,
		redis: redis,
		sf:    sf,
		rpc:   rpc,
	}
	return svc
}
