package rpc

import (
	"context"
	api "github.com/LingeringAutumn/Yijie/app/gateway/model/api/user"
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
