package usecase

import (
	"context"
	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
)

func (uc *userUseCase) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {}
