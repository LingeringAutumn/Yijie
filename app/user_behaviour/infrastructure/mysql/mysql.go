package mysql

import (
	"context"
	"errors"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/model"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/repository"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"gorm.io/gorm"
)

type userBehaviourDB struct {
	client *gorm.DB
}

func NewUserBehaviourDB(client *gorm.DB) repository.UserBehaviourDB {
	return &userBehaviourDB{client: client}
}

func (db *userBehaviourDB) LikeVideoDB(ctx context.Context, userID int64, videoID int64, isLike bool) error {
	var existing model.VideoLike

	// 1. 查询是否已有点赞记录
	err := db.client.WithContext(ctx).
		Table(model.VideoLike{}.TableName()).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query video_likes: %v", err)
	}

	if isLike {
		// === 点赞逻辑 ===
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 新增点赞记录
			newLike := model.VideoLike{
				UserID:  userID,
				VideoID: videoID,
				IsLiked: true,
			}
			if err := db.client.WithContext(ctx).
				Table(model.VideoLike{}.TableName()).
				Create(&newLike).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to insert video_like: %v", err)
			}
			// 点赞数 +1
			if err := db.client.WithContext(ctx).
				Table(constants.VideoStatsTableName).
				Where("video_id = ?", videoID).
				Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to increment like count: %v", err)
			}
		} else if !existing.IsLiked {
			// 记录存在但状态为取消，改为 true
			if err := db.client.WithContext(ctx).
				Table(model.VideoLike{}.TableName()).
				Where("user_id = ? AND video_id = ?", userID, videoID).
				Update("is_liked", true).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update video_like: %v", err)
			}
			// 点赞数 +1
			if err := db.client.WithContext(ctx).
				Table(constants.VideoStatsTableName).
				Where("video_id = ?", videoID).
				Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to increment like count: %v", err)
			}
		}
	} else {
		// === 取消点赞逻辑 ===
		if !errors.Is(err, gorm.ErrRecordNotFound) && existing.IsLiked {
			// 设置 is_liked = false
			if err := db.client.WithContext(ctx).
				Table(model.VideoLike{}.TableName()).
				Where("user_id = ? AND video_id = ?", userID, videoID).
				Update("is_liked", false).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to cancel like: %v", err)
			}
			// 点赞数 -1
			if err := db.client.WithContext(ctx).
				Table(constants.VideoStatsTableName).
				Where("video_id = ?", videoID).
				Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to decrement like count: %v", err)
			}
		}
	}

	return nil
}
