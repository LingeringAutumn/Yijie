package service

import (
	"github.com/LingeringAutumn/Yijie/app/user/domain/repository"
)

type UserService struct {
	db    repository.UserDB
	redis repository.UserRedis
	rpc   repository.UserRPC
}

func NewUserService(db repository.UserDB, redis repository.UserRedis, rpc repository.UserRPC) *UserService {
	if db == nil {
		panic("userService`s db should not be nil")
	}
	if redis == nil {
		panic("userService`s redis should not be nil")
	}
	if rpc == nil {
		panic("userService`s rpc should not be nil")
	}
	svc := &UserService{
		db:    db,
		redis: redis,
		rpc:   rpc,
	}
	return svc
}
