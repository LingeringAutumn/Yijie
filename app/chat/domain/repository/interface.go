package repository

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/LingeringAutumn/Yijie/app/chat/domain/model"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, msg *model.Message) error
	GetMessagesByChatID(ctx context.Context, chatID int64, limit, offset int) ([]*model.Message, error)
	ChatHandler(ctx context.Context, c *app.RequestContext)
}
