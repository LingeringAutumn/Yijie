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

// Components 封装 Handler 与 Service，用于依赖注入
type Components struct {
	Handler video.VideoService
	Service *service.VideoService
}

// InjectComponents 构建视频服务所有依赖并返回组件集合
func InjectComponents() *Components {
	// 初始化数据库存储
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}

	// 初始化 Redis 客户端
	redisClient, err := client.InitRedis(constants.RedisDBUser)
	if err != nil {
		panic(err)
	}
	redisRepo := redis.NewVideoRedis(redisClient)

	// 初始化雪花接口
	sf, err := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
	if err != nil {
		panic(err)
	}

	db := mysql.NewVideoDB(gormDB)
	svc := service.NewVideoService(db, redisRepo, sf)
	uc := usecase.NewVideoUseCase(db, redisRepo, sf, svc)
	handler := rpc.NewVideoHandler(uc)

	return &Components{
		Handler: handler,
		Service: svc,
	}
}
