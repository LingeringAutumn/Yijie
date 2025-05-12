package service

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/LingeringAutumn/Yijie/pkg/logger"

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

// TrendVideo Redis 热度排行榜 + 数据缓存
func (svc *VideoService) TrendVideo(ctx context.Context, pageNum, pageSize int64) ([]*model.VideoProfile, error) {
	start := (pageNum - 1) * pageSize
	end := start + pageSize - 1

	// 1. 从 Redis ZSet 获取热度排序前 N 的 video_id
	idsWithScores, err := svc.redis.GetHotRankRange(ctx, start, end)
	if err != nil {
		return nil, fmt.Errorf("redis hot_rank fetch failed: %w", err)
	}

	var results []*model.VideoProfile
	for _, item := range idsWithScores {
		videoIDStr, ok := item.Member.(string)
		if !ok {
			continue // 跳过非法 ID
		}
		videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
		if err != nil {
			continue // 跳过非法 ID
		}

		profile, err := svc.redis.GetVideoRedis(ctx, videoID)
		if err != nil {
			// Redis 未命中，查 DB
			profile, err = svc.db.GetVideoDB(ctx, videoID)
			if err != nil {
				continue
			}
			// 回写缓存
			_ = svc.redis.SetVideoRedis(ctx, profile)
		}
		profile.HotScore = item.Score // 从 ZSet 中带出的分数
		profile.Views, _ = svc.GetViews(ctx, videoID)
		results = append(results, profile)
	}

	return results, nil
}

// SearchVideo 按关键词模糊搜索，Redis 无法预缓存，仅搜索结果缓存
func (svc *VideoService) SearchVideo(ctx context.Context, keyword string, tags []string, pageNum, pageSize int64) ([]*model.VideoProfile, error) {
	cacheKey := fmt.Sprintf("video:search:%s:%d:%d", keyword, pageNum, pageSize)
	cached, err := svc.redis.GetSearchCache(ctx, cacheKey)
	if err == nil && cached != nil {
		return cached, nil
	}

	// SQL 模糊搜索
	dbResult, err := svc.db.SearchVideo(ctx, keyword, tags, pageNum, pageSize)
	if err != nil {
		return nil, err
	}

	// 写入缓存，设置过期时间 5 分钟
	_ = svc.redis.SetSearchCache(ctx, cacheKey, dbResult, 5*time.Minute)

	return dbResult, nil
}

func (svc *VideoService) IncrViews(ctx context.Context, videoId int64) (int64, error) {
	return svc.redis.IncrViews(ctx, videoId)
}

func (svc *VideoService) GetViews(ctx context.Context, videoId int64) (int64, error) {
	return svc.redis.GetViews(ctx, videoId)
}

func (svc *VideoService) UpdateHotRank(ctx context.Context, videoId int64, score float64) error {
	return svc.redis.UpdateHotRank(ctx, videoId, score)
}

func (svc *VideoService) UpdateHotScore(ctx context.Context, videoId int64, score float64) error {
	return svc.db.UpdateHotScore(ctx, videoId, score)
}

func (svc *VideoService) SyncViewsToDB(ctx context.Context) error {
	keys, err := svc.redis.ScanViewKeys(ctx)
	if err != nil {
		return fmt.Errorf("scan redis keys failed: %w", err)
	}

	for _, key := range keys {
		videoIDStr := strings.TrimPrefix(key, "video:views:")
		videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
		if err != nil {
			continue // 跳过非法 key
		}

		views, err := svc.redis.GetViews(ctx, videoID)
		if err != nil {
			continue // 跳过读取失败
		}

		if err := svc.db.UpdateViews(ctx, videoID, views); err != nil {
			return err
		}
	}

	return nil
}

// StartBackgroundTasks 启动播放量相关的定时写入与优雅退出任务
func (svc *VideoService) StartBackgroundTasks() {
	// 启动定时任务：每隔 5 分钟将 Redis 播放量同步写入 MySQL
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := svc.SyncViewsToDB(context.Background()); err != nil {
					logger.Errorf("periodic sync views failed: %v", err)
				}
			}
		}
	}()

	// 启动退出监听：服务关闭前执行一次 Redis → MySQL 同步
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		logger.Infof("termination signal received, flushing views to database")
		if err := svc.SyncViewsToDB(context.Background()); err != nil {
			logger.Errorf("flush before shutdown failed: %v", err)
		} else {
			logger.Infof("flush before shutdown success")
		}
		os.Exit(0)
	}()
}
