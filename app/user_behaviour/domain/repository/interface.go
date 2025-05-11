package repository

import "context"

type UserBehaviourDB interface {
	LikeVideoDB(ctx context.Context, userID int64, videoID int64, isLike bool) (err error)
}
