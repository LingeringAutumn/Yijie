package usecase

import (
	"context"
	"fmt"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/model"
)

func (uc *userBehaviourUseCase) LikeVideo(ctx context.Context, userBehaviour *model.VideoLike) (err error) {
	err = uc.svc.LikeVideo(ctx, userBehaviour.VideoID, userBehaviour.UserID, userBehaviour.IsLiked)
	if err != nil {
		return fmt.Errorf("usscase:like video failed: %w", err)
	}
	return err
}
