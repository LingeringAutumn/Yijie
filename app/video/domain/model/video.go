package model

import "time"

// Video 定义了用于存取数据库的核心视频结构
type Video struct {
	VideoID         int64      `json:"video_id" gorm:"primaryKey;column:video_id"`      // 视频ID
	UserID          int64      `json:"user_id" gorm:"column:user_id"`                   // 作者用户ID
	Title           string     `json:"title" gorm:"column:title"`                       // 视频标题
	Description     string     `json:"description" gorm:"column:description"`           // 视频描述
	CoverURL        string     `json:"cover_url" gorm:"column:cover_url"`               // 封面图URL
	VideoURL        string     `json:"video_url" gorm:"column:video_url"`               // 视频播放URL
	DurationSeconds int64      `json:"duration_seconds" gorm:"column:duration_seconds"` // 视频时长（单位：秒）
	Status          string     `json:"status" gorm:"column:status"`                     // 状态：published/deleted/draft
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at"`             // 发布时间
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at"`             // 更新时间
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`   // 逻辑删除时间，可空
}

func (Video) TableName() string {
	return "videos"
}

// VideoStat 定义了视频的统计信息表（单独维护）
type VideoStat struct {
	VideoID   int64     `json:"video_id" gorm:"primaryKey;column:video_id"`
	Views     int64     `json:"views" gorm:"column:views"`
	Likes     int64     `json:"likes" gorm:"column:likes"`
	Comments  int64     `json:"comments" gorm:"column:comments"`
	HotScore  float64   `json:"hot_score" gorm:"column:hot_score"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (VideoStat) TableName() string {
	return "video_stats"
}

// VideoProfile 聚合了视频主信息 + 统计信息，供展示用
type VideoProfile struct {
	VideoID         int64   `json:"video_id"`
	UserID          int64   `json:"user_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	CoverURL        string  `json:"cover_url"`
	VideoURL        string  `json:"video_url"`
	DurationSeconds int64   `json:"duration_seconds"`
	Views           int64   `json:"views"`
	Likes           int64   `json:"likes"`
	Comments        int64   `json:"comments"`
	HotScore        float64 `json:"hot_score"`
	CreatedAt       int64   `json:"created_at"` // 用时间戳方便前端
}
