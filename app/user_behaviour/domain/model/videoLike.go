package model

import "time"

// VideoLike 表示用户对视频的点赞记录
type VideoLike struct {
	UserID    int64     `json:"user_id" gorm:"primaryKey;column:user_id"`     // 用户ID
	VideoID   int64     `json:"video_id" gorm:"primaryKey;column:video_id"`   // 视频ID
	IsLiked   bool      `json:"is_liked" gorm:"column:is_liked;default:true"` // 是否点赞
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`          // 点赞时间
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`          // 更新时间
}

// TableName 返回该模型对应的数据库表名
func (VideoLike) TableName() string {
	return "video_likes"
}
