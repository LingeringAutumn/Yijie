package service

import (
	"context"
	"fmt"
)

func (svc *UserBehaviourService) LikeVideo(ctx context.Context, userID int64, videoID int64, isLike bool) error {
	// 1. 写入数据库
	if err := svc.db.LikeVideoDB(ctx, userID, videoID, isLike); err != nil {
		return fmt.Errorf("domain:like video failed: %w", err)
	}

	// 2. 启动异步任务处理 Redis + 热度值更新
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic recovered in LikeVideo async update: %v\n", r)
			}
		}()

		// 2.1 缓存点赞状态
		if isLike {
			if _, err := svc.redis.IncrLikes(context.Background(), videoID); err != nil {
				fmt.Printf("failed to incr likes in Redis: %v\n", err)
			}
			_ = svc.redis.SetLikeStatus(context.Background(), userID, videoID, true)
		} else {
			if _, err := svc.redis.DecrLikes(context.Background(), videoID); err != nil {
				fmt.Printf("failed to decr likes in Redis: %v\n", err)
			}
			_ = svc.redis.SetLikeStatus(context.Background(), userID, videoID, false)
		}

		// 2.2 跨模块 RPC：通知 video 模块更新热度值
		err := svc.rpc.UserBehaviourUpdateVideoHot(context.Background(), videoID)
		if err != nil {
			fmt.Printf("rpc UpdateVideoHot failed: %v\n", err)
		}
	}()
	return nil
}
