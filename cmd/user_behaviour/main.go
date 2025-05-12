package main

import (
	"log"
	"net"

	"github.com/LingeringAutumn/Yijie/app/user_behaviour"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user_behaviour/likeservice"

	"github.com/LingeringAutumn/Yijie/config"
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

var serviceName = constants.UserBehaviourServiceName

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
}

func main() {
	logger.Infof("starting user_behaviour service")
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("UserBehaviour: new etcd registry failed, err: %v", err)
	}

	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("UserBehaviour: get available port failed, err: %v", err)
	}
	log.Printf("UserBehaviour main running !!!!")
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("UserBehaviour: resolve tcp addr failed, err: %v", err)
	}
	err = utils.InitMinioClient(config.Minio.Endpoint, config.Minio.AccessKey, config.Minio.SecretKey)
	if err != nil {
		logger.Fatalf("UserBehaviour: new minio client failed, err: %v", err)
	}

	svr := likeservice.NewServer(
		// 注入依赖
		user_behaviour.InjectUserBehaviourHandler(),
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
	logger.Infof("UserBehaviour: server listening at %s", addr)

	if err = svr.Run(); err != nil {
		logger.Fatalf("user_behaviour: run server failed, err: %v", err)
	}
}
