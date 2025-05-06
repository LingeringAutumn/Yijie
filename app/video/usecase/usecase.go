package usecase

import (
	"context"
	"github.com/LingeringAutumn/Yijie/app/video/domain/model"
	"github.com/LingeringAutumn/Yijie/app/video/domain/repository"
	"github.com/LingeringAutumn/Yijie/app/video/domain/service"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

type VideoUseCase interface {
	SubmitVideo(ctx context.Context, video *model.Video, videoData []byte) (videoId int64, videoUrl string, err error)
	GetVideo(ctx context.Context, videoId int64) (*model.VideoProfile, error)
	SearchVideo(ctx context.Context, keyword string, tags []string, pageNum int64, pageSize int64) ([]*model.VideoProfile, error)
	TrendVideo(ctx context.Context, pageNum int64, pageSize int64) ([]*model.VideoProfile, error)
}

type videoUseCase struct {
	db    repository.VideoDB
	redis repository.VideoRedis
	sf    *utils.Snowflake
	svc   *service.VideoService
}

func NewVideoUseCase(db repository.VideoDB, redis repository.VideoRedis, sf *utils.Snowflake, svc *service.VideoService) VideoUseCase {
	return &videoUseCase{
		db:    db,
		redis: redis,
		sf:    sf,
		svc:   svc,
	}
}
