package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

func (svc *UserBehaviourService) LikeVideo(ctx context.Context, userID int64, videoID int64, isLike bool) error {
	// 1. 写入数据库
	if err := svc.db.LikeVideoDB(ctx, userID, videoID, isLike); err != nil {
		return fmt.Errorf("domain:like video failed: %w", err)
	}

	// 2. 异步更新 Redis 点赞状态、点赞数、热度值
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic recovered in LikeVideo cache update: %v\n", r)
			}
		}()

		// 2.1 点赞缓存 +1 / -1
		if isLike {
			_, _ = svc.redis.IncrLikes(ctx, videoID)
			_ = svc.redis.SetLikeStatus(ctx, userID, videoID, true)
		} else {
			_, _ = svc.redis.DecrLikes(ctx, videoID)
			_ = svc.redis.SetLikeStatus(ctx, userID, videoID, false)
		}

		// 2.2 获取播放量 + 点赞数
		views, _ := svc.redis.GetViews(ctx, videoID)
		likes, _ := svc.redis.GetLikes(ctx, videoID)

		// 2.3 获取视频创建时间（精确做法是从 Redis 缓存中取，如果你已有缓存）
		createdAt := time.Now()

		// 2.4 计算热度并更新
		hotScore := utils.ComputeHotScore(views, likes, createdAt)
		_ = svc.redis.UpdateHotRank(ctx, videoID, hotScore)
		_ = svc.db.UpdateHotScore(ctx, videoID, hotScore)
	}()

	return nil
}
