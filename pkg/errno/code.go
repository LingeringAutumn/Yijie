package errno

const (
	// SuccessCode For microservices
	SuccessCode = 10000
	SuccessMsg  = "success"
)

// 参数错误
const (
	ParamVerifyErrorCode   = 20000 + iota // 参数校验失败
	ParamMissingErrorCode                 // 参数缺失
	ParamMissingHeaderCode                // 请求头缺失
	ParamInvalidHeaderCode                // 请求头无效
)

// 鉴权错误
const (
	AuthInvalidCode         = 30000 + iota // 鉴权失败
	AuthPraiseKeyFailedCode                // token生成失败
	AuthSignKeyFailedCode
	AuthAccessExpiredCode       // 访问令牌过期
	AuthRefreshExpiredCode      // 刷新令牌过期
	AuthNoTokenCode             // 没有 token
	AuthNoOperatePermissionCode // 没有操作权限
	AuthMissingTokenCode        // 缺少 token
	IllegalOperatorCode         // 不合格的操作(比如传入 payment status时传入了一个不存在的 status)
)

// 内部错误,服务级别的错误
const (
	InternalServiceErrorCode  = 50000 + iota // 内部服务错误
	InternalDatabaseErrorCode                // 数据库错误
	InternalRedisErrorCode                   // Redis错误
	InternalNetworkErrorCode                 // 网络错误
	InternalKafkaErrorCode                   // kafka 错误
	InternalRPCErrorCode
)

// config错误
const (
	RedisConnectFailed = 60000 + iota
	RedisKeyNotExist
)
