package video

import (
	"github.com/LingeringAutumn/Yijie/app/video/controllers/rpc"
	"github.com/LingeringAutumn/Yijie/app/video/domain/service"
	"github.com/LingeringAutumn/Yijie/app/video/infrastructure/mysql"
	"github.com/LingeringAutumn/Yijie/app/video/infrastructure/redis"
	"github.com/LingeringAutumn/Yijie/app/video/usecase"
	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/kitex_gen/video"
	"github.com/LingeringAutumn/Yijie/pkg/base/client"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

// InjectUserHandler 用于依赖注入
// 从这个文件的位置就可以看出来极其特殊, 独立于架构之外, 服务于业务
func InjectVideoHandler() video.VideoService {
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
	redisRepo := redis.NewVideoRedis(redisClient)

	// 初始化雪花接口
	sf, err := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
	if err != nil {
		panic(err)
	}

	db := mysql.NewVideoDB(gormDB)
	svc := service.NewVideoService(db, redisRepo, sf)
	uc := usecase.NewVideoUseCase(db, redisRepo, sf, svc)

	return rpc.NewVideoHandler(uc)
}
