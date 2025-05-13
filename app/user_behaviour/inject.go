package user_behaviour

import (
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/controllers/rpc"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/domain/service"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/infrastructure/mysql"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/infrastructure/redis"
	"github.com/LingeringAutumn/Yijie/app/user_behaviour/usecase"
	"github.com/LingeringAutumn/Yijie/kitex_gen/user_behaviour"
	"github.com/LingeringAutumn/Yijie/pkg/base/client"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
)

// InjectUserBehaviourHandler 用于依赖注入
// 从这个文件的位置就可以看出来极其特殊, 独立于架构之外, 服务于业务
func InjectUserBehaviourHandler() user_behaviour.LikeService {
	// 初始化数据库存储
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}

	// 初始化 Redis 客户端
	// 初始化 Redis，使用指定的 Redis DB
	redisClient, err := client.InitRedis(constants.RedisDBUser)
	if err != nil {
		panic(err)
	}
	// 封装 Redis 存储对象
	redisRepo := redis.NewUserBehaviourRedis(redisClient)
	/*
		// 初始化雪花接口
		sf, err := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
		if err != nil {
			panic(err)
		}*/

	db := mysql.NewUserBehaviourDB(gormDB)
	svc := service.NewUserBehaviourService(db, redisRepo)
	uc := usecase.NewUserBehaviourUseCase(db, redisRepo, svc)

	return rpc.NewUserBehaviourHandler(uc)
}
