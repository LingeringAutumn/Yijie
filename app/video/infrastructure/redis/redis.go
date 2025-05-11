package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/LingeringAutumn/Yijie/app/video/domain/model"
	"github.com/LingeringAutumn/Yijie/app/video/domain/repository"
	"github.com/LingeringAutumn/Yijie/pkg/constants"

	"github.com/redis/go-redis/v9"
)

type videoRedis struct {
	client *redis.Client
}

func NewVideoRedis(client *redis.Client) repository.VideoRedis {
	cli := videoRedis{client: client}
	return &cli
}
func (v *videoRedis) GetVideoRedis(ctx context.Context, videoId int64) (*model.VideoProfile, error) {
	key := fmt.Sprintf("video:%d", videoId)

	// 从 Redis 获取数据
	result, err := v.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("video not found in cache: %w", err)
		}
		return nil, fmt.Errorf("redis get failed: %w", err)
	}

	// 反序列化 JSON
	var video model.VideoProfile
	err = json.Unmarshal([]byte(result), &video)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	// ✅ 空结构校验
	if video.VideoID == 0 {
		return nil, fmt.Errorf("redis cache hit empty video")
	}

	return &video, nil
}

func (v *videoRedis) SetVideoRedis(ctx context.Context, videoProfile *model.VideoProfile) error {
	if videoProfile == nil || videoProfile.VideoID == 0 {
		return fmt.Errorf("refuse to cache empty video profile")
	}

	key := fmt.Sprintf("video:%d", videoProfile.VideoID)

	videoJSON, err := json.Marshal(videoProfile)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	err = v.client.Set(ctx, key, videoJSON, constants.RedisHalfHour*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}

	return nil
}
