package controllers

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
)

var (
	connMap  = make(map[int64]*websocket.Conn)
	mu       sync.Mutex
	Upgrader = websocket.HertzUpgrader{
		CheckOrigin: func(c *app.RequestContext) bool {
			return true
		},
	}
)

// 这个函数必须存在，并且名字必须是大写开头才能被导出
func ChatHandler(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.String(400, "invalid user_id")
		return
	}

	Upgrader.Upgrade(c, func(conn *websocket.Conn) {
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
		}
	})
}
