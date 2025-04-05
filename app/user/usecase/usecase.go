package usecase

import (
	"context"
	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
	"github.com/LingeringAutumn/Yijie/app/user/domain/repository"
	"github.com/LingeringAutumn/Yijie/app/user/domain/service"
)

// UserUseCase 接口应该不应该定义在 domain 中，这属于 use case 层
type UserUseCase interface {
	RegisterUser(ctx context.Context, user *model.User) (uid int64, err error)
	LoginUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUserProfile(ctx context.Context, user *model.UserProfileRequest) (*model.UserProfileResponse, error)
	GetUserProfile(ctx context.Context, uid int64) (*model.UserProfileResponse, error)
}

type userUseCase struct {
	db    repository.UserDB
	redis repository.UserRedis
	rpc   repository.UserRPC
	svc   *service.UserService
}

func NewUserUseCase(db repository.UserDB, svc *service.UserService, redis repository.UserRedis, rpc repository.UserRPC) UserUseCase {
	return &userUseCase{
		db:    db,
		svc:   svc,
		redis: redis,
		rpc:   rpc,
	}
}
