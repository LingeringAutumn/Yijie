package constants

import "time"

const (
	MuxConnection    = 1                     // (RPC) 最大连接数
	StreamBufferSize = 1024                  // (RPC) 流请求 Buffer 尺寸
	MaxQPS           = 100                   // (RPC) 最大 QPS
	RPCTimeout       = 3 * time.Second       // (RPC) RPC请求超时时间
	ConnectTimeout   = 50 * time.Millisecond // (RPC) 连接超时时间

	KitexClientEndpointInfoFormat = "%s-client" // serviceName-client
)
