package rpc

import (
	"context"

	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/repository"
	videorpc "github.com/LingeringAutumn/Yijie/kitex_gen/video"
	"github.com/LingeringAutumn/Yijie/kitex_gen/video/videoservice"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

type userBehaviourRPC struct {
	video videoservice.Client
}

func NewUserBehaviourRPC(video videoservice.Client) repository.UserBehaviourRPC {
	return &userBehaviourRPC{video: video}
}

func (rpc *userBehaviourRPC) UserBehaviourUpdateVideoHot(ctx context.Context, videoId int64) (err error) {
	videoRpcReq := &videorpc.VideoHotUpdateRequest{
		VideoId: videoId,
	}
	resp, err := rpc.video.UpdateVideoHot(ctx, videoRpcReq)
	if err = utils.ProcessRpcError("user_behaviour.UpdateVideoHot", resp, err); err != nil {
		return err
	}
	return nil
}
