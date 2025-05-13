package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	dmodel "github.com/LingeringAutumn/Yijie/app/video/domain/model"
)

type VideoDB interface {
	StoreVideo(ctx context.Context, video *dmodel.Video) error
	GetVideoDB(ctx context.Context, videoId int64) (*dmodel.VideoProfile, error)
	SearchVideo(ctx context.Context, keyword string, tags []string, num int64, size int64) ([]*dmodel.VideoProfile, error)
	TrendVideo(ctx context.Context, num int64, size int64) ([]*dmodel.VideoProfile, error)
	StoreVideoStats(ctx context.Context, stat *dmodel.VideoStat) error
	UpdateViews(ctx context.Context, videoID int64, views int64) error
	UpdateHotScore(ctx context.Context, videoID int64, score float64) error
}

type VideoRedis interface {
	GetVideoRedis(ctx context.Context, videoId int64) (*dmodel.VideoProfile, error)
	SetVideoRedis(ctx context.Context, videoProfile *dmodel.VideoProfile) error
	IncrViews(ctx context.Context, videoID int64) (int64, error)
	GetViews(ctx context.Context, videoID int64) (int64, error)
	ScanViewKeys(ctx context.Context) ([]string, error)
	UpdateHotRank(ctx context.Context, videoID int64, hotScore float64) error
	GetHotRankRange(ctx context.Context, start, end int64) ([]redis.Z, error)
	SetSearchCache(ctx context.Context, key string, data []*dmodel.VideoProfile, ttl time.Duration) error
	GetSearchCache(ctx context.Context, key string) ([]*dmodel.VideoProfile, error)
}

type VideoRPC interface{}
