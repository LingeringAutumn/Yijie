package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/LingeringAutumn/Yijie/pkg/utils"

	"github.com/LingeringAutumn/Yijie/app/gateway/router"
	"github.com/LingeringAutumn/Yijie/app/gateway/rpc"
	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/logger"
)

// 定义服务名称，从 constants 包中获取
var serviceName = constants.GatewayServiceName

// init 函数会在 main 函数之前自动执行，用于初始化配置、日志和 RPC 等
func init() {
	// 初始化配置，传入服务名称
	config.Init(serviceName)
	// 初始化日志，传入服务名称和从配置中获取的日志级别
	logger.Init(serviceName, config.GetLoggerLevel())
	// 初始化 RPC 相关配置
	rpc.Init()
}

func main() {
	// 调用 utils 包中的 GetAvailablePort 函数，获取一个可用的端口
	// 如果获取过程中出现错误，使用 logger.Fatalf 打印错误信息并终止程序
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	// listenAddr := "0.0.0.0:8888"
	logger.Infof("available port:%s", listenAddr)

	// 创建一个新的 Hertz 服务器实例
	// server.WithHostPorts(listenAddr)：设置服务器监听的地址和端口
	// server.WithHandleMethodNotAllowed(true)：允许处理 HTTP 方法不允许的情况
	// server.WithMaxRequestBodySize(constants.ServerMaxRequestBodySize)：设置服务器允许的最大请求体大小，该值从 constants 包中获取
	h := server.New(
		server.WithHostPorts(listenAddr),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(constants.ServerMaxRequestBodySize),
	)
	// 调用 router 包中的 GeneratedRegister 函数，将路由规则注册到 Hertz 服务器上
	router.GeneratedRegister(h)
	// 启动 Hertz 服务器，开始监听客户端请求
	h.Spin()
}
