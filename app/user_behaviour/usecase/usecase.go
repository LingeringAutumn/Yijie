package usecase

import (
	"context"

	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/model"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/repository"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/service"
)

type UserBehaviourUseCase interface {
	LikeVideo(ctx context.Context, userBehaviour *model.VideoLike) (err error)
}
type userBehaviourUseCase struct {
	db    repository.UserBehaviourDB
	redis repository.UserBehaviourRedis
	svc   *service.UserBehaviourService
	rpc   repository.UserBehaviourRPC
}

func NewUserBehaviourUseCase(db repository.UserBehaviourDB, redis repository.UserBehaviourRedis, svc *service.UserBehaviourService, rpc repository.UserBehaviourRPC) UserBehaviourUseCase {
	return &userBehaviourUseCase{
		db:    db,
		redis: redis,
		svc:   svc,
		rpc:   rpc,
	}
}
