package repository

import (
	"context"
	dmodel "github.com/LingeringAutumn/Yijie/app/video/domain/model"
)

type VideoDB interface {
	StoreVideo(ctx context.Context, video *dmodel.Video) error
}

type VideoRedis interface {
}

type VideoRPC interface{}
