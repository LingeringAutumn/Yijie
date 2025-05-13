package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/repository"
)

type userBehaviourRedis struct {
	client *redis.Client
}

func NewUserBehaviourRedis(client *redis.Client) repository.UserBehaviourRedis {
	cli := userBehaviourRedis{client: client}
	return &cli
}

func (r *userBehaviourRedis) SetLikeStatus(ctx context.Context, userID, videoID int64, isLike bool) error {
	key := fmt.Sprintf("video:like:%d:%d", userID, videoID)
	val := "0"
	if isLike {
		val = "1"
	}
	return r.client.Set(ctx, key, val, 24*time.Hour).Err()
}

func (r *userBehaviourRedis) GetLikeStatus(ctx context.Context, userID, videoID int64) (bool, error) {
	key := fmt.Sprintf("video:like:%d:%d", userID, videoID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return val == "1", nil
}

func (r *userBehaviourRedis) IncrLikes(ctx context.Context, videoID int64) (int64, error) {
	return r.client.Incr(ctx, fmt.Sprintf("video:likes:%d", videoID)).Result()
}

func (r *userBehaviourRedis) DecrLikes(ctx context.Context, videoID int64) (int64, error) {
	return r.client.Decr(ctx, fmt.Sprintf("video:likes:%d", videoID)).Result()
}

func (r *userBehaviourRedis) GetViews(ctx context.Context, videoID int64) (int64, error) {
	key := fmt.Sprintf("video:views:%d", videoID)
	val, err := r.client.Get(ctx, key).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return val, nil
}

func (r *userBehaviourRedis) GetLikes(ctx context.Context, videoID int64) (int64, error) {
	key := fmt.Sprintf("video:likes:%d", videoID)
	val, err := r.client.Get(ctx, key).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return val, nil
}

func (r *userBehaviourRedis) UpdateHotRank(ctx context.Context, videoID int64, hotScore float64) error {
	return r.client.ZAdd(ctx, "video:hot_rank", redis.Z{
		Score:  hotScore,
		Member: videoID,
	}).Err()
}
