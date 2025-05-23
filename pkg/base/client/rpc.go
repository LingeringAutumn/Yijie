package client

import (
	"errors"
	"fmt"

	"github.com/LingeringAutumn/Yijie/kitex_gen/video/videoservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user/userservice"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user_behaviour/likeservice"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
)

// 通用的RPC客户端初始化函数
func initRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	// 初始化Etcd Resolver
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("initRPCClient etcd.NewEtcdResolver failed: %w", err)
	}
	// 初始化具体的RPC客户端
	client, err := newClientFunc(serviceName,
		client.WithResolver(r),
		client.WithMuxConnection(constants.MuxConnection),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: fmt.Sprintf(constants.KitexClientEndpointInfoFormat, serviceName)}),
	)
	if err != nil {
		return nil, fmt.Errorf("initRPCClient NewClient failed: %w", err)
	}
	return &client, nil
}

func InitUserRPC() (*userservice.Client, error) {
	return initRPCClient(constants.UserServiceName, userservice.NewClient)
}

func InitVideoRPC() (*videoservice.Client, error) {
	return initRPCClient(constants.VideoServiceName, videoservice.NewClient)
}

func InitUserBehaviourRPC() (*likeservice.Client, error) {
	return initRPCClient(constants.UserBehaviourServiceName, likeservice.NewClient)
}
