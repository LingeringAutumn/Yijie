package constants

const (
    DefaultVideoCoverUrl = ""
    RedisMinute          = 60
    RedisHalfHour        = 1800
)

const (
    DecayFactor          float64 = 3600 * 6 // 每 6 小时衰减一分
    DefaultHotScoreDelta         = -2
    HotRankKey                   = "video:hot_rank"
)

const (
	DecayFactor float64 = 3600 * 6 // 每 6 小时衰减一分
	HotRankKey          = "video:hot_rank"
)
