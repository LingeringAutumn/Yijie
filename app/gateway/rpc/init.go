package rpc

import "github.com/LingeringAutumn/Yijie/kitex_gen/user/userservice"

var userClient userservice.Client

func Init() {
	InitUserRPC()
}
