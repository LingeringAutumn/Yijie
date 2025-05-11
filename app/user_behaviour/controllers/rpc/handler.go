package rpc

import (
	"context"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/usecase"

	"github.com/LingeringAutumn/Yijie/kitex_gen/user_behaviour"
	"github.com/LingeringAutumn/Yijie/pkg/base"
)

type UserBehaviourHandler struct {
	useCase usecase.UserBehaviourUseCase
}

func NewUserBehaviourHandler(useCase usecase.UserBehaviourUseCase) *UserBehaviourHandler {
	return &UserBehaviourHandler{useCase: useCase}
}

func (handler *UserBehaviourHandler) LikeVideo(ctx context.Context, req *user_behaviour.VideoLikeRequest) (resp *user_behaviour.VideoLikeResponse, err error) {
	resp = new(user_behaviour.VideoLikeResponse)
	err = handler.useCase.LikeVideo(ctx, req.VideoId, req.UserId, req.IsLike)
	if err != nil {
		resp.BaseResp = base.BuildBaseResp(err)
		return
	}
	resp.BaseResp = base.BuildBaseResp(err)
	return
}
