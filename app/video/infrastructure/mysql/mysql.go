package mysql

import (
	"context"
	dmodel "github.com/LingeringAutumn/Yijie/app/video/domain/model"
	"github.com/LingeringAutumn/Yijie/app/video/domain/repository"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"gorm.io/gorm"
)

type videoDB struct {
	client *gorm.DB
}

func NewVideoDB(client *gorm.DB) repository.VideoDB {
	return &videoDB{client: client}
}

func (db *videoDB) StoreVideo(ctx context.Context, video *dmodel.Video) error {
	if err := db.client.WithContext(ctx).Table(constants.VideoTableName).Create(&video).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to store video: %v", err)
	}
	return nil
}
func (db *videoDB) GetVideoDB(ctx context.Context, videoId int64) (*dmodel.VideoProfile, error) {
	var video dmodel.Video
	var stat dmodel.VideoStat

	// 1. 查询视频主表（videos）
	err := db.client.WithContext(ctx).
		Table(constants.VideoTableName).
		Where("video_id = ?", videoId).
		First(&video).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "query videos failed: %v", err)
	}

	// 2. 查询视频统计表（video_stats）
	err = db.client.WithContext(ctx).
		Table(constants.VideoStatsTableName).
		Where("video_id = ?", videoId).
		First(&stat).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "query video_stats failed: %v", err)
	}

	// 3. 聚合数据构造 VideoProfile 返回结构
	profile := &dmodel.VideoProfile{
		VideoID:         video.VideoID,
		UserID:          video.UserID,
		Title:           video.Title,
		Description:     video.Description,
		CoverURL:        video.CoverURL,
		VideoURL:        video.VideoURL,
		DurationSeconds: video.DurationSeconds,
		Views:           stat.Views,
		Likes:           stat.Likes,
		Comments:        stat.Comments,
		HotScore:        stat.HotScore,
		CreatedAt:       video.CreatedAt.Unix(), // 转换为时间戳（秒）
	}

	return profile, nil
}
