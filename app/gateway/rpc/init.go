package rpc

import (
	"github.com/LingeringAutumn/Yijie/kitex_gen/user/userservice"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user_behaviour/likeservice"
	"github.com/LingeringAutumn/Yijie/kitex_gen/video/videoservice"
)

var (
	userClient  userservice.Client
	videoClient videoservice.Client
	likeClient  likeservice.Client
)

func Init() {
	InitUserRPC()
	InitVideoRPC()
	InitUserBehaviourRPC()
}
