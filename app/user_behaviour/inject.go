package user_behaviour

import (
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/controllers/rpc"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/service"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/infrastructure/mysql"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/infrastructure/redis"
	videoRpcPkg "github.com/LingeringAutumn/Yijie/app/user_behaviour/infrastructure/rpc"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/usecase"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user_behaviour"
	"github.com/LingeringAutumn/Yijie/pkg/base/client"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/logger"

)

// InjectUserBehaviourHandler 用于依赖注入
// 从这个文件的位置就可以看出来极其特殊, 独立于架构之外, 服务于业务
func InjectUserBehaviourHandler() user_behaviour.LikeService {
	// 初始化数据库存储
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	// 初始化 Redis，使用指定的 Redis DB
	redisClient, err := client.InitRedis(constants.RedisDBUser)
	if err != nil {
		panic(err)
	}
	// 封装 Redis 存储对象
	redisRepo := redis.NewUserBehaviourRedis(redisClient)
	c, err := client.InitVideoRPC()
	if err != nil {
		logger.Fatalf("api.rpc.video InitVideoRPC failed, err is %v", err)
	}
	videoRpc := videoRpcPkg.NewUserBehaviourRPC(*c)
	db := mysql.NewUserBehaviourDB(gormDB)
	svc := service.NewUserBehaviourService(db, redisRepo, videoRpc)
	uc := usecase.NewUserBehaviourUseCase(db, redisRepo, svc, videoRpc)

	return rpc.NewUserBehaviourHandler(uc)
}
