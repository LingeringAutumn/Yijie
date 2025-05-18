package service

import (
	"context"

	"github.com/LingeringAutumn/Yijie/app/chat/domain/model"
	"github.com/LingeringAutumn/Yijie/app/chat/domain/repository"
)

type ChatService struct {
	repo repository.MessageRepository
}

func NewChatService(repo repository.MessageRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) SaveMessage(ctx context.Context, msg *model.Message) error {
	return s.repo.SaveMessage(ctx, msg)
}
