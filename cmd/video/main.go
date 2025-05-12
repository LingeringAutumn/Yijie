package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/LingeringAutumn/Yijie/app/video/domain/service"

	"github.com/LingeringAutumn/Yijie/app/video"
	"github.com/LingeringAutumn/Yijie/kitex_gen/video/videoservice"

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

var serviceName = constants.VideoServiceName

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
}

func waitForExitAndFlush(svc *service.VideoService) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logger.Infof("Video: Received an exit signal, writing the view count immediately")
	if err := svc.SyncViewsToDB(context.Background()); err != nil {
		logger.Errorf("Video: Failed to write the view count: %v", err)
	} else {
		logger.Infof("Video: Successfully wrote the view count, exiting")
	}
	os.Exit(0)
}

func main() {
	logger.Infof("starting video service")
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Video: new etcd registry failed, err: %v", err)
	}

	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Video: get available port failed, err: %v", err)
	}
	log.Printf("video main running !!!!")
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Video: resolve tcp addr failed, err: %v", err)
	}
	err = utils.InitMinioClient(config.Minio.Endpoint, config.Minio.AccessKey, config.Minio.SecretKey)
	if err != nil {
		logger.Fatalf("Video: new minio client failed, err: %v", err)
	}

	components := video.InjectComponents()

	components.Service.StartBackgroundTasks()

	svr := videoservice.NewServer(
		// 注入依赖
		components.Handler,
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
	logger.Infof("Video: server listening at %s", addr)

	if err = svr.Run(); err != nil {
		logger.Fatalf("Video: run server failed, err: %v", err)
	}
}
