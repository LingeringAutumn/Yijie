package router

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/websocket"
)

var (
	connMap  = make(map[int64]*websocket.Conn)
	mu       sync.Mutex
	upgrader = websocket.HertzUpgrader{
		// 允许所有来源连接
		CheckOrigin: func(c *app.RequestContext) bool {
			return true
		},
	}
)

// RegisterChatRoutes 注册聊天 WebSocket 路由
func RegisterChatRoutes(h *server.Hertz) {
	h.GET("/ws/:user_id", func(ctx context.Context, c *app.RequestContext) {
		userIDStr := c.Param("user_id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.String(400, "invalid user_id")
			return
		}

		upgrader.Upgrade(c, func(conn *websocket.Conn) {
			defer conn.Close()

			mu.Lock()
			connMap[userID] = conn
			mu.Unlock()

			log.Printf("[Chat] User %d connected.\n", userID)

			for {
				var msg string
				err := conn.ReadJSON(&msg)
				if err != nil {
					log.Printf("[Chat] User %d disconnected: %v\n", userID, err)
					mu.Lock()
					delete(connMap, userID)
					mu.Unlock()
					break
				}

				log.Printf("[Chat] Received from %d: %s\n", userID, msg)
				// 👉 下一步：转发到对方、持久化数据库等
			}
		})
	})
}
