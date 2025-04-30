package rpc

import (
	"golang.org/x/net/context"

	"github.com/LingeringAutumn/Yijie/pkg/logger"

	"github.com/LingeringAutumn/Yijie/app/user/controllers/rpc/pack"
	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
	"github.com/LingeringAutumn/Yijie/app/user/usecase"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user"
	"github.com/LingeringAutumn/Yijie/pkg/base"
)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

func (handler *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (r *user.RegisterResponse, err error) {
	r = new(user.RegisterResponse)
	logger.Info("new user register")
	u := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	var uid int64
	if uid, err = handler.useCase.RegisterUser(ctx, u); err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	logger.Info("user register success")
	r.Base = base.BuildBaseResp(err)
	r.UserID = uid
	return
}

func (handler *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (r *user.LoginResponse, err error) {
	r = new(user.LoginResponse)
	u := &model.User{
		Username: req.Username,
		Password: req.Password,
	}
	ans, err := handler.useCase.LoginUser(ctx, u)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(err)
	r.User = pack.BuildUser(ans)
	return
}

func (handler *UserHandler) UpdateProfile(ctx context.Context, req *user.UpdateUserProfileRequest) (r *user.UpdateUserProfileResponse, err error) {
	r = new(user.UpdateUserProfileResponse)
	u := &model.UserProfileRequest{
		Username: req.UserProfileReq.Username,
		Email:    req.UserProfileReq.Email,
		Phone:    req.UserProfileReq.Phone,
		Bio:      req.UserProfileReq.Bio,
	}
	avatar := req.Avatar
	// TODO team,membershipLevel和point应该在UpdateUserProfile函数里处理，这些东西不应该由前端传过来
	userProfile, err := handler.useCase.UpdateUserProfile(ctx, u, avatar)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(err)
	r.UserProfileResp = pack.BuildUpdateUserProfileResponse(userProfile)
	return
}

func (handler *UserHandler) GetProfile(ctx context.Context, req *user.GetUserProfileRequest) (r *user.GetUserProfileResponse, err error) {
	r = new(user.GetUserProfileResponse)
	userProfile, err := handler.useCase.GetUserProfile(ctx, req.Uid)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(err)
	r.UserProfileResp = pack.BuildGetUserProfileResponse(userProfile)
	return
}
