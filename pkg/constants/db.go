package constants

import "time"

const (
	UserTableName       = "users"
	ImageTableName      = "images"
	VideoTableName      = "videos"
	VideoStatsTableName = "video_stats"
)

const (
	// MaxConnections 表示数据库的最大连接数。
	// 这是数据库连接池允许同时存在的最大连接数量。
	// 当并发请求较多时，连接数可能会达到这个上限，
	// 超过该上限的请求可能需要等待其他连接释放。
	MaxConnections = 1000 // (DB) 最大连接数

	// MaxIdleConns 表示数据库的最大空闲连接数。
	// 连接池会保持一定数量的空闲连接，以便快速响应新的请求。
	// 当空闲连接数超过这个值时，多余的空闲连接会被关闭。
	MaxIdleConns = 10 // (DB) 最大空闲连接数

	// ConnMaxLifetime 表示数据库连接的最大可复用时间。
	// 一个连接在被创建后，如果使用时间超过这个值，
	// 连接池会将其关闭，避免连接长时间占用资源。
	ConnMaxLifetime = 10 * time.Second // (DB) 最大可复用时间

	// ConnMaxIdleTime 表示数据库连接最长保持空闲状态的时间。
	// 如果一个连接在空闲状态下的时间超过这个值，
	// 连接池会将其关闭，以释放资源。
	ConnMaxIdleTime = 5 * time.Minute // (DB) 最长保持空闲状态时间
)
