package main

import (
	"github.com/LingeringAutumn/Yijie/app/chat/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	"log"
)

func main() {
	h := server.Default(server.WithHostPorts(":8889")) // 自定义端口避免与其他模块冲突

	// 注册聊天 WebSocket 路由
	router.RegisterChatRoutes(h)
	log.Println("Chat service started at :8889")
	// 启动服务
	h.Spin()
}
