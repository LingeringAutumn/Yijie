package constants

import "time"

// RedisSlowQuery 是 Redis 默认的慢查询时间，单位为毫秒，用于日志记录。
// 当 Redis 执行某个命令的时间超过这个值时，可认为该查询为慢查询。
const (
	RedisSlowQuery = 10 // ms redis默认的慢查询时间，适用于 logger
)

// Redis Key and Expire Time
// 以下常量定义了 Redis 中不同键的过期时间和存储数量相关信息
const (
	// RedisNXExpireTime 是使用 Redis 的 NX（仅在键不存在时设置）操作时的过期时间，为 3 秒。
	RedisNXExpireTime = 3 * time.Second
	// RedisMaxLockRetryTime 是获取 Redis 锁时的最大重试时间，为 400 毫秒。
	RedisMaxLockRetryTime = 400 * time.Millisecond
	// RedisRetryStopTime 是获取 Redis 锁时，重试停止的时间间隔，为 100 毫秒。
	RedisRetryStopTime = 100 * time.Millisecond
)

// Redis DB Name
// 以下常量定义了 Redis 不同数据库的编号
const (
	// RedisDBUser 表示用于存储User相关数据的 Redis 数据库编号为 0。
	RedisDBUser = 0

	// RedSyncDBId 是 RedSync（分布式锁库）使用的 Redis 数据库编号为 0。
	RedSyncDBId = 0
)

// Redis Connection Pool Configuration
// 以下常量定义了 Redis 连接池的相关配置参数
const (
	// RedisPoolSize 是 Redis 连接池的最大连接数，设置为 50。
	RedisPoolSize = 50 // 最大连接数
	// RedisMinIdleConnections 是 Redis 连接池的最小空闲连接数，为 10。
	RedisMinIdleConnections = 10 // 最小空闲连接数
	// RedisDialTimeout 是连接 Redis 时的超时时间，设置为 5 秒。
	RedisDialTimeout = 5 * time.Second // 连接超时时间
)

// 以下常量用于表示 Redis 的健康状态、检查超时时间和检查间隔时间
const (
	// RedisUnHealthy 表示 Redis 处于不健康状态，值为 false。
	RedisUnHealthy = false
	// RedisHealthy 表示 Redis 处于健康状态，值为 true。
	RedisHealthy = true
	// RedisCheckoutTimeOut 是检查 Redis 健康状态的超时时间，为 2 秒。
	RedisCheckoutTimeOut = 2 * time.Second
	// RedisCheckoutInterval 是检查 Redis 健康状态的时间间隔，为 5 秒。
	RedisCheckoutInterval = 5 * time.Second
)

// PingTime TODO 这个可能要删改吧
const (
	PingTime = 2
)
