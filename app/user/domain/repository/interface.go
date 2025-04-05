package repository

import (
	"context"
	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
)

type UserDB interface {
	IsUserExist(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, u *model.User) (int64, error)
	GetUserById(ctx context.Context, uid int64) (*model.User, error)
	GetUserProfileInfoById(ctx context.Context, uid int64) (*model.UserProfileResponse, error)
	StoreUserAvatar(ctx context.Context, image *model.Image) error
	StoreUserProfile(ctx context.Context, userProfileRequest *model.UserProfileRequest, uid int64, image *model.Image) (*model.UserProfileResponse, error)
}

type UserRedis interface {
}

type UserRPC interface {
}
