package service

import "github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/repository"

type UserBehaviourService struct {
	db    repository.UserBehaviourDB
	redis repository.UserBehaviourRedis
	rpc   repository.UserBehaviourRPC
}

func NewUserBehaviourService(db repository.UserBehaviourDB, redis repository.UserBehaviourRedis, rpc repository.UserBehaviourRPC) *UserBehaviourService {
	if db == nil {
		panic("userBehaviourService`s db should not be nil")
	}
	if redis == nil {
		panic("userBehaviourService`s redis should not be nil")
	}
	if rpc == nil {
		panic("userBehaviourService`s rpc should not be nil")
	}
	svc := &UserBehaviourService{
		db:    db,
		redis: redis,
		rpc:   rpc,
	}
	return svc
}
