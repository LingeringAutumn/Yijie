package main

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github/LinegringAutumn/Yijie/pkg/constants"
)

var serviceName = constants.GatewayServiceName

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	rpc.Init()
}

func main() {
	// get available port from config set
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("get available port failed, err: %v", err)
	}

	h := server.New(
		server.WithHostPorts(listenAddr),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(constants.ServerMaxRequestBodySize),
	)

	h.Use(
		mw.RecoveryMW(), // recovery
		mw.CorsMW(),     // cors
		mw.GzipMW(),     // gzip
	)

	router.GeneratedRegister(h)
	h.Spin()
}
