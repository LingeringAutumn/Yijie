package main

import (
	"net"

	"github.com/LingeringAutumn/Yijie/app/user"
	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user/userservice"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/logger"
	"github.com/LingeringAutumn/Yijie/pkg/middleware"
	"github.com/LingeringAutumn/Yijie/pkg/utils"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var serviceName = constants.UserServiceName

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("User: new etcd registry failed, err: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("User: get available port failed, err: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("User: resolve tcp addr failed, err: %v", err)
	}
	err = utils.InitMinioClient(config.Minio.Addr, config.Minio.AccessKeyID, config.Minio.AccessKey)
	if err != nil {
		logger.Fatalf("User: new minio client failed, err: %v", err)
	}

	svr := userservice.NewServer(
		// 注入依赖
		user.InjectUserHandler(),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: serviceName,
		}),
		server.WithMuxTransport(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{
			MaxConnections: constants.MaxConnections,
			MaxQPS:         constants.MaxQPS,
		}),

		server.WithMiddleware(middleware.Respond()),
	)
	if err = svr.Run(); err != nil {
		logger.Fatalf("User: run server failed, err: %v", err)
	}
}
