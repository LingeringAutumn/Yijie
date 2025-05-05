package service

import (
	"context"
	"fmt"
	userData "github.com/LingeringAutumn/Yijie/pkg/base/context"
	"math"
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
