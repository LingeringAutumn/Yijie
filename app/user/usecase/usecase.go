package usecase

import "context"

// UserUseCase 接口应该不应该定义在 domain 中，这属于 use case 层
type UserUseCase interface {
	RegisterUser(ctx context.Context, user *model.User) (uid int64, err error)
	Login(ctx context.Context, user *model.User) (*model.User, error)
}
