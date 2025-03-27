package rpc

import (
	"context"
	api "github.com/LingeringAutumn/Yijie/app/gateway/model/api/user"
	"github.com/LingeringAutumn/Yijie/app/gateway/model/model"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user"
	"github.com/LingeringAutumn/Yijie/pkg/base/client"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
	"github.com/bytedance/gopkg/util/logger"
)

func InitUserRPC() {
	c, err := client.InitUserRPC()
	if err != nil {
		logger.Fatal("api.rpc.user InitUserRPC failed, err is %v", err)
	}
	userClient = *c
}

func RegisterRPC(ctx context.Context, req *user.RegisterRequest) (response *api.RegisterResponse, err error) {
	resp, err := userClient.Register(ctx, req)
	// 这里的 err 是属于 RPC 间调用的错误，例如 network error
	// 而业务错误则是封装在 resp.base 当中的
	if err != nil {
		logger.Fatal("RegisterRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	// 用中间件去判断resp.Base里是否有错误
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithError(err)
	}
	response = &api.RegisterResponse{UID: resp.UserID}
	return response, nil
}

func LoginRPC(ctx context.Context, req *user.LoginRequest) (response *api.LoginResponse, err error) {
	resp, err := userClient.Login(ctx, req)
	if err != nil {
		logger.Fatal("LoginRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithError(err)
	}

	response = &api.LoginResponse{
		User: &model.UserInfo{
			UserId: resp.User.UserId,
			Name:   resp.User.Name,
		},
	}
	return response, nil
}

func UpdateUserProfileRPC(ctx context.Context, req *user.UpdateUserProfileRequest) (response *api.UpdateUserProfileResponse, err error) {
	resp, err := userClient.UpdateUserProfile(ctx, req)
	if err != nil {
		logger.Fatal("UpdateUserProfileRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithError(err)
	}
	response = &api.UpdateUserProfileResponse{
		UserProfile: &model.UserProfile{
			Username:        resp.UserProfile.Username,
			Email:           resp.UserProfile.Email,
			Phone:           resp.UserProfile.Phone,
			Avatar:          resp.UserProfile.Avatar,
			Bio:             resp.UserProfile.Bio,
			MembershipLevel: resp.UserProfile.MembershipLevel,
			Point:           resp.UserProfile.Point,
			Team:            resp.UserProfile.Team,
		},
	}
	return response, nil
}

func GetUserProfileRPC(ctx context.Context, req *user.GetUserProfileRequest) (response *api.GetUserProfileResponse, err error) {
	resp, err := userClient.GetUserProfile(ctx, req)
	if err != nil {
		logger.Fatal("GetUserProfile: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithError(err)
	}
	response = &api.GetUserProfileResponse{
		UserProfile: &model.UserProfile{
			Username:        resp.UserProfile.Username,
			Email:           resp.UserProfile.Email,
			Phone:           resp.UserProfile.Phone,
			Avatar:          resp.UserProfile.Avatar,
			Bio:             resp.UserProfile.Bio,
			MembershipLevel: resp.UserProfile.MembershipLevel,
			Point:           resp.UserProfile.Point,
			Team:            resp.UserProfile.Team,
		},
	}
	return response, nil
}
