package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/LingeringAutumn/Yijie/app/video/domain/model"
	"github.com/LingeringAutumn/Yijie/pkg/constants"

	"github.com/redis/go-redis/v9"
	"time"
)

type videoRedis struct {
	client *redis.Client
}

func (v *videoRedis) GetVideoRedis(ctx context.Context, videoId int64) (*model.VideoProfile, error) {
	key := fmt.Sprintf("video:%d", videoId)

	var video model.VideoProfile
	err := v.client.Get(ctx, key).Scan(&video)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("video not found in cache: %w", err)
		}
		return nil, fmt.Errorf("redis get failed: %w", err)
	}
	return &video, nil
}

func (v *videoRedis) SetVideoRedis(ctx context.Context, videoProfile *model.VideoProfile) error {
	key := fmt.Sprintf("video:%d", videoProfile.VideoID)

	// 设置缓存，超时时间30分钟
	err := v.client.Set(ctx, key, videoProfile, constants.RedisHalfHour*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}
	return nil
}
