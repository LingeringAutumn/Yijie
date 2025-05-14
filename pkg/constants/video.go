package constants

const (
	DefaultVideoCoverUrl = ""
	RedisMinute          = 60
	RedisHalfHour        = 1800
)

const (
	DecayFactor          float64 = 3600 * 6 // 每 6 小时衰减一分
	DefaultHotScoreDelta         = -1
	HotRankKey                   = "video:hot_rank"
)
