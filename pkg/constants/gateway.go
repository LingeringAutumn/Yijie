package constants

import "time"

const (
	// 请求体最大体积
	ServerMaxRequestBodySize = 1 << 31

	CorsMaxAge = 12 * time.Hour

	SentinelThreshold        = 100
	SentinelStatIntervalInMs = 1000
	LoginDataKey             = "loginData"
)
