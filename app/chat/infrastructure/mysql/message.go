package mysql

import (
	"context"
	"fmt"

	"github.com/LingeringAutumn/Yijie/app/chat/domain/model"
	"github.com/LingeringAutumn/Yijie/app/chat/domain/repository"
	"gorm.io/gorm"
)

type messageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) repository.MessageRepository {
	return &messageRepo{db: db}
}

func (r *messageRepo) SaveMessage(ctx context.Context, msg *model.Message) error {
	if err := r.db.WithContext(ctx).Table("messages").Create(msg).Error; err != nil {
		return fmt.Errorf("save message failed: %w", err)
	}
	return nil
}

func (r *messageRepo) GetMessagesByChatID(ctx context.Context, chatID int64, limit, offset int) ([]*model.Message, error) {
	var msgs []*model.Message
	if err := r.db.WithContext(ctx).
		Table("messages").
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}
