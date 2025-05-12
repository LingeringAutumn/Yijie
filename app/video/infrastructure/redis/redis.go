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

// IncrViews 播放量 +1，返回新值
func (v *videoRedis) IncrViews(ctx context.Context, videoID int64) (int64, error) {
	key := fmt.Sprintf("video:views:%d", videoID)
	return v.client.Incr(ctx, key).Result()
}

// GetViews 获取当前播放量
func (v *videoRedis) GetViews(ctx context.Context, videoID int64) (int64, error) {
	key := fmt.Sprintf("video:views:%d", videoID)
	val, err := v.client.Get(ctx, key).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("get views from redis failed: %w", err)
	}
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return val, nil
}

// ScanViewKeys 扫描所有 video:views:* 的 Redis 键
func (v *videoRedis) ScanViewKeys(ctx context.Context) ([]string, error) {
	iter := v.client.Scan(ctx, 0, "video:views:*", 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func (v *videoRedis) UpdateHotRank(ctx context.Context, videoID int64, hotScore float64) error {
	return v.client.ZAdd(ctx, constants.HotRankKey, redis.Z{
		Score:  hotScore,
		Member: videoID,
	}).Err()
}

// GetHotRankRange 热榜读取
func (v *videoRedis) GetHotRankRange(ctx context.Context, start, end int64) ([]redis.Z, error) {
	return v.client.ZRevRangeWithScores(ctx, "video:hot_rank", start, end).Result()
}

// GetSearchCache 搜索缓存读取
func (v *videoRedis) GetSearchCache(ctx context.Context, key string) ([]*model.VideoProfile, error) {
	val, err := v.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if val == "" {
		return nil, nil
	}
	var result []*model.VideoProfile
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SetSearchCache 搜索缓存写入
func (v *videoRedis) SetSearchCache(ctx context.Context, key string, data []*model.VideoProfile, ttl time.Duration) error {
	if len(data) == 0 {
		return nil // 空结果不缓存
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return v.client.Set(ctx, key, bytes, ttl).Err()
}
