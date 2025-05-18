package usecase

import (
	"context"

	"github.com/LingeringAutumn/Yijie/app/chat/domain/model"
	"github.com/LingeringAutumn/Yijie/app/chat/domain/service"
)

type ChatUseCase struct {
	svc *service.ChatService
}

func NewChatUseCase(svc *service.ChatService) *ChatUseCase {
	return &ChatUseCase{svc: svc}
}

func (uc *ChatUseCase) HandleMessage(ctx context.Context, msg *model.Message) error {
	return uc.svc.SaveMessage(ctx, msg)
}
