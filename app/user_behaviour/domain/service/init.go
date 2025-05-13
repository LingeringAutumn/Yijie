package service

import "github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/repository"

type UserBehaviourService struct {
	db    repository.UserBehaviourDB
	redis repository.UserBehaviourRedis
}

func NewUserBehaviourService(db repository.UserBehaviourDB, redis repository.UserBehaviourRedis) *UserBehaviourService {
	if db == nil {
		panic("userBehaviourService`s db should not be nil")
	}
	if redis == nil {
		panic("userBehaviourService`s redis should not be nil")
	}
	svc := &UserBehaviourService{
		db:    db,
		redis: redis,
	}
	return svc
}
