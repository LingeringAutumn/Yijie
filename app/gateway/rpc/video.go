package rpc

import (
	"context"

	"github.com/bytedance/gopkg/util/logger"

	"github.com/LingeringAutumn/Yijie/kitex_gen/video"
	"github.com/LingeringAutumn/Yijie/pkg/base/client"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

func InitVideoRPC() {
	c, err := client.InitVideoRPC()
	if err != nil {
		logger.Fatal("api.rpc.video InitVideoRPC failed, err is %v", err)
	}
	videoClient = *c
}

func SubmitVideoRPC(ctx context.Context, req *video.VideoSubmissionRequest) (*video.VideoSubmissionResponse, error) {
	resp, err := videoClient.SubmitVideo(ctx, req)
	if err != nil {
		logger.Errorf("SubmitVideoRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.BaseResp) {
		return nil, errno.InternalServiceError.WithMessage(resp.BaseResp.Msg)
	}
	return resp, nil
}

func GetVideoRPC(ctx context.Context, req *video.VideoDetailRequest) (*video.VideoDetailResponse, error) {
	resp, err := videoClient.GetVideo(ctx, req)
	if err != nil {
		logger.Errorf("SubmitVideoRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.BaseResp) {
		return nil, errno.InternalServiceError.WithMessage(resp.BaseResp.Msg)
	}
	return resp, nil
}

func SearchVideoRPC(ctx context.Context, req *video.VideoSearchRequest) (*video.VideoSearchResponse, error) {
	resp, err := videoClient.SearchVideo(ctx, req)
	if err != nil {
		logger.Errorf("SearchVideoRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.BaseResp) {
		return nil, errno.InternalServiceError.WithMessage(resp.BaseResp.Msg)
	}
	return resp, nil
}

func TrendVideoRPC(ctx context.Context, req *video.VideoTrendingRequest) (*video.VideoTrendingResponse, error) {
	resp, err := videoClient.TrendVideo(ctx, req)
	if err != nil {
		logger.Errorf("TrendVideoRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.BaseResp) {
		return nil, errno.InternalServiceError.WithMessage(resp.BaseResp.Msg)
	}
	return resp, nil
}
