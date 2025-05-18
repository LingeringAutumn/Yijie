package repository

import "context"

type UserBehaviourDB interface {
	LikeVideoDB(ctx context.Context, userID int64, videoID int64, isLike bool) (err error)


type UserBehaviourRedis interface {
	SetLikeStatus(ctx context.Context, userID, videoID int64, isLike bool) error
	GetLikeStatus(ctx context.Context, userID, videoID int64) (bool, error)
	IncrLikes(ctx context.Context, videoID int64) (int64, error)
	DecrLikes(ctx context.Context, videoID int64) (int64, error)
	GetLikes(ctx context.Context, videoID int64) (int64, error)
	UpdateHotRank(ctx context.Context, videoID int64, hotScore float64) error
	GetViews(ctx context.Context, videoID int64) (int64, error)
}

type UserBehaviourRPC interface {
	UserBehaviourUpdateVideoHot(ctx context.Context, videoId int64) (err error)
}
