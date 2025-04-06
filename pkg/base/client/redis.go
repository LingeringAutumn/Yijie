package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"github.com/LingeringAutumn/Yijie/pkg/logger"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

// InitRedis 函数用于初始化 Redis 客户端。
// 它接收一个整数参数 db，表示要使用的 Redis 数据库编号。
// 函数返回一个指向 redis.Client 的指针和一个错误对象。
func InitRedis(db int) (*redis.Client, error) {
	// 检查 Redis 配置是否为空。
	// 如果配置为空，说明没有正确配置 Redis，返回一个错误。
	if config.Redis == nil {
		return nil, errors.New("redis config is nil")
	}

	// 创建一个新的 Redis 客户端实例。
	// 使用 config.Redis 中的地址、密码和传入的数据库编号。
	// 同时设置连接池的大小、最小空闲连接数和连接超时时间，这些参数从 constants 包中获取。
	client := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Password,
		DB:           db,
		PoolSize:     constants.RedisPoolSize,           // 连接池大小
		MinIdleConns: constants.RedisMinIdleConnections, // 最小空闲连接数
		DialTimeout:  constants.RedisDialTimeout,        // 连接超时
	})

	// 添加日志 Hook。
	// 获取 Redis 日志记录器，并将其设置为 Redis 客户端的日志记录器。
	// 同时将日志记录器添加为客户端的 Hook，以便记录客户端的操作。
	l := logger.GetRedisLogger()
	redis.SetLogger(l)
	client.AddHook(l)

	// 使用超时的 Ping 检查 Redis 连接是否正常。
	// 创建一个带有超时的上下文，超时时间从 constants 包中获取。
	// 使用该上下文执行 Ping 命令，如果出现错误，返回一个自定义的错误。
	ctx, cancel := context.WithTimeout(context.Background(), constants.PingTime*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("client.NewRedisClient: ping redis failed: %v", err))
	}

	// 如果一切正常，返回创建的 Redis 客户端和 nil 错误。
	return client, nil
}

// InitRedSync 函数用于初始化 RedSync 实例。
// RedSync 是一个用于在 Redis 上实现分布式锁的库。
// 它接收一个指向 redis.Client 的指针作为参数。
// 函数返回一个指向 redsync.Redsync 的指针。
func InitRedSync(client *redis.Client) *redsync.Redsync {
	// 创建一个基于 Redis 客户端的连接池。
	// 使用 goredis.NewPool 函数将 Redis 客户端包装成 RedSync 所需的连接池。
	pool := goredis.NewPool(client)
	// 使用创建的连接池初始化 RedSync 实例。
	// RedSync 使用该连接池与 Redis 进行交互，以实现分布式锁。
	rs := redsync.New(pool)
	// 返回初始化好的 RedSync 实例。
	return rs
}
