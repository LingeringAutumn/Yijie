package constants

import (
	"time"
)

const (
	TypeAccessToken  = 0
	TypeRefreshToken = 1

	AccessTokenTTL  = time.Hour * 24 * 7  // Access Token 有效期7天
	RefreshTokenTTL = time.Hour * 24 * 30 // Refresh Token 有效期30天
	// Issuer TODO 这个地方是不是要换成我的名字
	Issuer = "yijie" // token 颁发者

	AuthHeader         = "Authorization" // 获取 Token 时的请求头
	AccessTokenHeader  = "Access-Token"  // 响应时的访问令牌头
	RefreshTokenHeader = "Refresh-Token" // 响应时的刷新令牌头
)
