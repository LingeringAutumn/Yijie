package repository

import (
	"context"
	dmodel "github.com/LingeringAutumn/Yijie/app/video/domain/model"
)

type VideoDB interface {
	StoreVideo(ctx context.Context, video *dmodel.Video) error
	GetVideoDB(ctx context.Context, videoId int64) (*dmodel.VideoProfile, error)
	SearchVideo(ctx context.Context, keyword string, tags []string, num int64, size int64) ([]*dmodel.VideoProfile, error)
	TrendVideo(ctx context.Context, num int64, size int64) ([]*dmodel.VideoProfile, error)
}

type VideoRedis interface {
	GetVideoRedis(ctx context.Context, videoId int64) (*dmodel.VideoProfile, error)
	SetVideoRedis(ctx context.Context, videoProfile *dmodel.VideoProfile) error
}

type VideoRPC interface{}
