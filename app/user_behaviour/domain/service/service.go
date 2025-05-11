package service

import (
	"context"
	"fmt"
	userData "github.com/LingeringAutumn/Yijie/pkg/base/context"
)

func (svc *UserBehaviourService) LikeVideo(ctx context.Context, userID int64, videoID int64, isLike bool) (err error) {
	err = svc.db.LikeVideoDB(ctx, userID, videoID, isLike)
	if err != nil {
		return fmt.Errorf("domain:like video failed: %w", err)
	}
	return err
}

func (svc *UserBehaviourService) GetUserId(ctx context.Context) (uid int64, err error) {
	if uid, err = userData.GetLoginData(ctx); err != nil {
		return 0, fmt.Errorf("get user id failed: %w", err)
	}
	return uid, err
}
