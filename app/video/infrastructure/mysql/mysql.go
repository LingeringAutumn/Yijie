package mysql

import (
	"context"
	"fmt"
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
	var video dmodel.Video    // 存储视频主表信息
	var stat dmodel.VideoStat // 存储视频统计信息

	// 查询视频主表（videos），根据 video_id 精确查找
	err := db.client.WithContext(ctx).
		Table(constants.VideoTableName).
		Where("video_id = ?", videoId).
		First(&video).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "query videos failed: %v", err)
	}

	// 查询视频统计表（video_stats），同样根据 video_id 精确查找
	err = db.client.WithContext(ctx).
		Table(constants.VideoStatsTableName).
		Where("video_id = ?", videoId).
		First(&stat).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "query video_stats failed: %v", err)
	}

	// 聚合主表 + 统计表信息，构造返回结果结构体
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

func (db *videoDB) SearchVideo(ctx context.Context, keyword string, tags []string, pageNum, pageSize int64) ([]*dmodel.VideoProfile, error) {
	var results []*dmodel.VideoProfile
	offset := int((pageNum - 1) * pageSize) // 分页偏移计算

	// 执行多表联查（videos + video_stats），查询符合关键词的视频
	err := db.client.WithContext(ctx).
		Table(fmt.Sprintf("%s AS v", constants.VideoTableName)). // 主表别名 v
		Select(`
			v.video_id, v.user_id, v.title, v.description, v.cover_url, v.video_url,
			v.duration_seconds, UNIX_TIMESTAMP(v.created_at) as created_at,
			vs.views, vs.likes, vs.comments, vs.hot_score
		`). // SELECT 字段来自主表和统计表
		Joins(fmt.Sprintf("LEFT JOIN %s AS vs ON v.video_id = vs.video_id", constants.VideoStatsTableName)). // 统计信息表联接
		Where("v.status = ?", "published"). // 只查询已发布视频
		Where("v.title LIKE ? OR v.description LIKE ?", "%"+keyword+"%", "%"+keyword+"%"). // 模糊搜索标题或描述
		Order("v.created_at DESC"). // 新发布的在前
		Offset(offset). // 分页偏移
		Limit(int(pageSize)). // 限制结果数量
		Scan(&results).Error // 结果扫描进结构体切片

	if err != nil {
		return nil, fmt.Errorf("search videos failed: %w", err)
	}

	return results, nil
}

func (db *videoDB) TrendVideo(ctx context.Context, pageNum, pageSize int64) ([]*dmodel.VideoProfile, error) {
	var results []*dmodel.VideoProfile
	offset := int((pageNum - 1) * pageSize) // 分页偏移计算

	// 查询热榜视频：按 hot_score 排序
	err := db.client.WithContext(ctx).
		Table(fmt.Sprintf("%s AS v", constants.VideoTableName)). // 主表别名 v
		Select(`
			v.video_id, v.user_id, v.title, v.description, v.cover_url, v.video_url,
			v.duration_seconds, UNIX_TIMESTAMP(v.created_at) as created_at,
			vs.views, vs.likes, vs.comments, vs.hot_score
		`). // SELECT 字段来自主表和统计表
		Joins(fmt.Sprintf("LEFT JOIN %s AS vs ON v.video_id = vs.video_id", constants.VideoStatsTableName)). // 联接统计表
		Where("v.status = ?", "published"). // 只显示已发布视频
		Order("vs.hot_score DESC, v.created_at DESC"). // 按热度排序，发布时间为次要排序
		Offset(offset).
		Limit(int(pageSize)).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("fetch trending videos failed: %w", err)
	}

	return results, nil
}
