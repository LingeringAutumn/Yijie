package rpc

import (
	"context"

	"github.com/LingeringAutumn/Yijie/app/video/controllers/rpc/pack"
	"github.com/LingeringAutumn/Yijie/app/video/domain/model"
	"github.com/LingeringAutumn/Yijie/app/video/usecase"
	"github.com/LingeringAutumn/Yijie/kitex_gen/video"
	"github.com/LingeringAutumn/Yijie/pkg/base"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
)

type VideoHandler struct {
	useCase usecase.VideoUseCase
}

func NewVideoHandler(useCase usecase.VideoUseCase) *VideoHandler {
	return &VideoHandler{useCase: useCase}
}

func (handler *VideoHandler) SubmitVideo(ctx context.Context, req *video.VideoSubmissionRequest) (resp *video.VideoSubmissionResponse, err error) {
	resp = new(video.VideoSubmissionResponse)
	var cover string
	if req.CoverUrl != nil {
		cover = *req.CoverUrl
	} else {
		cover = constants.DefaultVideoCoverUrl
	}
	v := &model.Video{
		Title:           req.Title,
		Description:     req.Description,
		DurationSeconds: req.DurationSeconds,
		CoverURL:        cover,
	}
	videoId, videoUrl, err := handler.useCase.SubmitVideo(ctx, v, req.Video)
	if err != nil {
		resp.BaseResp = base.BuildBaseResp(err)
		return
	}
	resp.VideoId = videoId
	resp.VideoUrl = videoUrl
	return
}

func (handler *VideoHandler) GetVideo(ctx context.Context, req *video.VideoDetailRequest) (resp *video.VideoDetailResponse, err error) {
	resp = new(video.VideoDetailResponse)
	videoProfile, err := handler.useCase.GetVideo(ctx, req.VideoId)
	if err != nil {
		resp.BaseResp = base.BuildBaseResp(err)
		return
	}
	resp.Video = pack.BuildVideo(videoProfile)
	return
}

func (handler *VideoHandler) SearchVideo(ctx context.Context, req *video.VideoSearchRequest) (resp *video.VideoSearchResponse, err error) {
	resp = new(video.VideoSearchResponse)
	videoList, err := handler.useCase.SearchVideo(ctx, req.Keyword, req.Tags, req.PageNum, req.PageSize)
	if err != nil {
		resp.BaseResp = base.BuildBaseResp(err)
		return
	}
	resp.Videos = pack.BuildVideoList(videoList)
	return
}

func (handler *VideoHandler) TrendVideo(ctx context.Context, req *video.VideoTrendingRequest) (resp *video.VideoTrendingResponse, err error) {
	resp = new(video.VideoTrendingResponse)
	videoList, err := handler.useCase.TrendVideo(ctx, req.PageNum, req.PageSize)
	if err != nil {
		resp.BaseResp = base.BuildBaseResp(err)
		return
	}
	resp.Videos = pack.BuildVideoList(videoList)
	return
}
