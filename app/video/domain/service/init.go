package service

import "github.com/LingeringAutumn/Yijie/app/video/domain/repository"

type VideoService struct {
	db    repository.VideoDB
	redis repository.VideoRedis
}

func NewVideoService(db repository.VideoDB, redis repository.VideoRedis) *VideoService {
	if db == nil {
		panic("videoService`s db should not be nil")
	}
	if redis == nil {
		panic("videoervice`s redis should not be nil")
	}
	svc := &VideoService{
		db:    db,
		redis: redis,
	}
	return svc
}
