package main

import (
	"github.com/LingeringAutumn/Yijie/app/gateway/mw"
	"github.com/LingeringAutumn/Yijie/app/gateway/router"
	"github.com/LingeringAutumn/Yijie/app/gateway/rpc"
	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/logger"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app/server"
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
		logger.Fatalf("get available port failed, err: %v", err)
	}

	// 创建一个遥测提供者，用于分布式追踪和监控
	// 传入服务名称和从配置中获取的收集器地址
	p := base.TelemetryProvider(serviceName, config.Otel.CollectorAddr)
	// 使用 defer 确保在 main 函数结束时关闭遥测提供者
	defer func() { logger.LogError(p.Shutdown(context.Background())) }()

	// 创建一个新的 Hertz 服务器实例
	// server.WithHostPorts(listenAddr)：设置服务器监听的地址和端口
	// server.WithHandleMethodNotAllowed(true)：允许处理 HTTP 方法不允许的情况
	// server.WithMaxRequestBodySize(constants.ServerMaxRequestBodySize)：设置服务器允许的最大请求体大小，该值从 constants 包中获取
	h := server.New(
		server.WithHostPorts(listenAddr),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(constants.ServerMaxRequestBodySize),
	)

	// 为 Hertz 服务器添加中间件
	// mw.RecoveryMW()：用于恢复 panic 错误，避免服务器崩溃
	// mw.CorsMW()：处理跨域资源共享（CORS）问题
	// mw.GzipMW()：对响应进行 Gzip 压缩，减少数据传输量
	// mw.SentinelMW()：使用 Sentinel 进行流量控制和熔断降级
	h.Use(
		mw.RecoveryMW(), // recovery
		mw.CorsMW(),     // cors
		mw.GzipMW(),     // gzip
		mw.SentinelMW(), // sentinel
	)

	// 调用 router 包中的 GeneratedRegister 函数，将路由规则注册到 Hertz 服务器上
	router.GeneratedRegister(h)
	// 启动 Hertz 服务器，开始监听客户端请求
	h.Spin()