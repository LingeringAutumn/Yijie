package rpc

import (
	"context"

	"github.com/bytedance/gopkg/util/logger"

	"github.com/LingeringAutumn/Yijie/kitex_gen/user_behaviour"
	"github.com/LingeringAutumn/Yijie/pkg/base/client"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

func InitUserBehaviourRPC() {
	c, err := client.InitUserBehaviourRPC()
	if err != nil {
		logger.Fatal("api.rpc.user_behaviour InitUserBehaviourRPC failed, err is %v", err)
	}
	likeClient = *c
}

func LikeVideoRPC(ctx context.Context, req *user_behaviour.VideoLikeRequest) (*user_behaviour.VideoLikeResponse, error) {
	resp, err := likeClient.LikeVideo(ctx, req)
	if err != nil {
		logger.Errorf("LikeVideoRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.BaseResp) {
		return nil, errno.InternalServiceError.WithMessage(resp.BaseResp.Msg)
	}
	return resp, nil
}
