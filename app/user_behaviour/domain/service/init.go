package service

import "github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/repository"

type UserBehaviourService struct {
	db repository.UserBehaviourDB
}

func NewUserBehaviourService(db repository.UserBehaviourDB) *UserBehaviourService {
	if db == nil {
		panic("userBehaviourService`s db should not be nil")
	}
	svc := &UserBehaviourService{
		db: db,
	}
	return svc
}
