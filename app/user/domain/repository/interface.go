package repository

import (
	"context"
	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
)

type UserDB interface {
	IsUserExist(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, u *model.User) (int64, error)
	GetUserById(ctx context.Context, uid int64) (*model.User, error)
	GetUserProfileInfoById(ctx context.Context, uid int64) (*model.UserProfile, error)
}

type UserRedis interface {
}

type UserRPC interface {
}
