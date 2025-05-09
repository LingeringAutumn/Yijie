package service

import (
	"context"
	"fmt"
	"math"

	"github.com/LingeringAutumn/Yijie/app/video/domain/model"
	userData "github.com/LingeringAutumn/Yijie/pkg/base/context"
)

func (svc *VideoService) GetUserId(ctx context.Context) (uid int64, err error) {
	if uid, err = userData.GetLoginData(ctx); err != nil {
		return 0, fmt.Errorf("get user id failed: %w", err)
	}
	return uid, err
}

func (svc *VideoService) GenerateVideoId() (VideoId int64, err error) {
	videoId, err := svc.sf.NextVal()
	if err != nil {
		return 0, fmt.Errorf("sf: failed to generate video ID: %w", err)
	}

	// 可选：检查是否超出 int64 范围（通常不会）
	if videoId > math.MaxInt64 {
		return 0, fmt.Errorf("sf: generated ID exceeds int64 limit")
	}

	return videoId, nil
}

func (svc *VideoService) StoreVideo(ctx context.Context, video *model.Video) (err error) {
	err = svc.db.StoreVideo(ctx, video)
	if err != nil {
		return fmt.Errorf("store video failed: %w", err)
	}
	return err
}

func (svc *VideoService) StoreVideoStats(ctx context.Context, stat *model.VideoStat) error {
	err := svc.db.StoreVideoStats(ctx, stat)
	if err != nil {
		return fmt.Errorf("store video stats failed: %w", err)
	}
	return nil
}

func (svc *VideoService) GetVideoDB(ctx context.Context, videoId int64) (*model.VideoProfile, error) {
	return svc.db.GetVideoDB(ctx, videoId)
}

func (svc *VideoService) GetVideoRedis(ctx context.Context, videoId int64) (*model.VideoProfile, error) {
	return svc.redis.GetVideoRedis(ctx, videoId)
}

func (svc *VideoService) SetVideoRedis(ctx context.Context, videoProfile *model.VideoProfile) error {
	return svc.redis.SetVideoRedis(ctx, videoProfile)
}

func (svc *VideoService) SearchVideo(ctx context.Context, keyword string, tags []string, pageNum int64, pageSize int64) ([]*model.VideoProfile, error) {
	return svc.db.SearchVideo(ctx, keyword, tags, pageNum, pageSize)
}

func (svc *VideoService) TrendVideo(ctx context.Context, pageNum int64, pageSize int64) ([]*model.VideoProfile, error) {
	return svc.db.TrendVideo(ctx, pageNum, pageSize)
}
