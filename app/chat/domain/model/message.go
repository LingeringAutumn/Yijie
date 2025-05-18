package model

import "time"

type Message struct {
    ID         int64     `json:"id"`
    ChatID     int64     `json:"chat_id"`
    SenderID   int64     `json:"sender_id"`
    ReceiverID int64     `json:"receiver_id"`
    Content    string    `json:"content"`
    Type       string    `json:"type"` // text/image/voice/...
    CreatedAt  time.Time `json:"created_at"`
}

type ChatMessage struct {
    SenderID   int64  `json:"sender_id"`
    ReceiverID int64  `json:"receiver_id"`
    Content    string `json:"content"`
    IsGroup    bool   `json:"is_group"` // 暂不实现群聊，后续扩展
}
