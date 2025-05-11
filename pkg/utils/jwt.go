package utils

// 鉴权相关组件
import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
)

// Claims 定义 JWT 令牌的声明结构体，包含令牌类型、用户 ID 以及标准声明
type Claims struct {
	Type               int64 `json:"type"`    // 令牌类型，用于区分不同类型的令牌
	UserID             int64 `json:"user_id"` // 用户 ID，标识令牌所属的用户
	jwt.StandardClaims       // 嵌入标准的 JWT 声明
}

// CreateToken 根据 token 类型和用户 ID 创建 token
// 参数 tokenType 为令牌类型，uid 为用户 ID
// 返回值为生成的令牌字符串和可能出现的错误
func CreateToken(tokenType int64, uid int64) (string, error) {
	// 检查配置中的服务器信息是否存在
	if config.Server == nil {
		return "", errno.NewErrNo(errno.AuthInvalidCode, "config server not found")
	}

	// 计算令牌的过期时间，调用 getTokenTTL 函数获取过期时长
	expireTime := time.Now().Add(getTokenTTL(tokenType))
	// 构建令牌声明
	claims := Claims{
		Type:   tokenType,
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间戳
			IssuedAt:  time.Now().Unix(), // 签发时间戳
			Issuer:    constants.Issuer,  // 签发者
		},
	}

	// 使用 EdDSA 签名方法创建 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	// 解析私钥
	key, err := parsePrivateKey(config.Server.Secret)
	if err != nil {
		return "", errno.Errorf(errno.AuthPraiseKeyFailedCode, "failed to parse private key: %v", err)
	}

	// 对令牌进行签名
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", errno.Errorf(errno.AuthSignKeyFailedCode, "failed to sign token: %v", err)
	}

	return signedToken, nil
}

// CreateAllToken 创建一对 token，第一个是 access token，第二个是 refresh token
// 参数 uid 为用户 ID
// 返回值为生成的访问令牌、刷新令牌和可能出现的错误
func CreateAllToken(uid int64) (string, string, error) {
	// 创建访问令牌
	accessToken, err := CreateToken(constants.TypeAccessToken, uid)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}
	// 创建刷新令牌
	refreshToken, err := CreateToken(constants.TypeRefreshToken, uid)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// CheckToken 检查 token 是否有效，并返回 token 类型和用户 ID
// 参数 token 为要检查的令牌字符串
// 返回值为令牌类型、用户 ID 和可能出现的错误
func CheckToken(token string) (int64, int64, error) {
	// 检查配置中的服务器信息是否存在
	if config.Server == nil {
		return 0, 0, errno.NewErrNo(errno.AuthInvalidCode, "server config not found")
	}
	// 检查令牌是否为空
	if token == "" {
		return -1, 0, errno.NewErrNo(errno.AuthMissingTokenCode, "token is empty")
	}

	// 解析未验证的令牌声明
	unverifiedClaims, err := parseUnverifiedClaims(token)
	if err != nil {
		return -1, 0, fmt.Errorf("failed to parse unverified claims: %w", err)
	}

	// 解析公钥
	secret, err := parsePublicKey(config.Server.PublicKey)
	if err != nil {
		return -1, 0, fmt.Errorf("failed to parse public key: %w", err)
	}

	// 验证令牌并获取验证后的声明
	verifiedClaims, err := verifyToken(token, secret)
	if err != nil {
		// 处理令牌验证错误
		tokenType, err := handleTokenError(err, unverifiedClaims.Type)
		return tokenType, 0, err
	}

	return verifiedClaims.Type, verifiedClaims.UserID, nil
}

// parsePrivateKey 解析 Ed25519 私钥
// 参数 key 为私钥的 PEM 格式字符串
// 返回值为解析后的私钥接口和可能出现的错误
// Create时传入的key为私钥，Verify时传入的key为公钥
// 下面的函数显然是不需要被外部访问的，所以首字母小写
func parsePrivateKey(key string) (interface{}, error) {
	// 从 PEM 格式的字符串中解析 Ed25519 私钥
	privateKey, err := jwt.ParseEdPrivateKeyFromPEM([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return privateKey, nil
}

// parsePublicKey 解析 Ed25519 公钥
// 参数 key 为公钥的 PEM 格式字符串
// 返回值为解析后的公钥接口和可能出现的错误
func parsePublicKey(key string) (interface{}, error) {
	// 从 PEM 格式的字符串中解析 Ed25519 公钥
	publicKey, err := jwt.ParseEdPublicKeyFromPEM([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return publicKey, nil
}

// parseUnverifiedClaims 解析未验证的 token claims
// 参数 token 为要解析的令牌字符串
// 返回值为解析后的声明结构体指针和可能出现的错误
func parseUnverifiedClaims(token string) (*Claims, error) {
	// 解析未验证的令牌
	tokenStruct, _, err := new(jwt.Parser).ParseUnverified(token, &Claims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse unverified token: %w", err)
	}

	// 类型断言，将解析后的声明转换为自定义的 Claims 结构体
	claims, ok := tokenStruct.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("failed to parse unverified claims:%w", err)
	}
	return claims, nil
}

// verifyToken 验证 token 并返回 claims
// 参数 token 为要验证的令牌字符串，key 为用于验证的公钥
// 返回值为验证后的声明结构体指针和可能出现的错误
func verifyToken(token string, key interface{}) (*Claims, error) {
	// 解析并验证令牌
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法是否为 Ed25519
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// 类型断言，检查解析后的令牌是否有效，并将声明转换为自定义的 Claims 结构体
	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("failed to parse token:%v", err)
}

// handleTokenError 处理 token 验证错误
// 参数 err 为验证过程中出现的错误，tokenType 为令牌类型
// 返回值为处理后的令牌类型和可能出现的错误
func handleTokenError(err error, tokenType int64) (int64, error) {
	// 声明一个指向 jwt.ValidationError 类型的指针 ve，用于后续存储类型断言后的错误对象
	var ve *jwt.ValidationError

	// 使用 errors.As 函数进行类型断言，检查传入的错误 err 是否为 jwt.ValidationError 类型
	// 如果是，则将 err 转换为 jwt.ValidationError 类型并赋值给 ve，同时返回 true
	// 如果不是，则 ve 保持为 nil，返回 false
	if errors.As(err, &ve) {
		// ve.Errors 是一个位掩码，用于表示具体的验证错误类型
		// jwt.ValidationErrorExpired 是一个预定义的常量，表示令牌过期错误
		// 使用按位与运算符 & 检查 ve.Errors 中是否包含 jwt.ValidationErrorExpired 标志
		// 如果结果不为 0，则说明存在令牌过期错误
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			// 检查令牌类型是否为访问令牌
			if tokenType == constants.TypeAccessToken {
				// 如果是访问令牌过期，返回 -1 作为处理后的令牌类型
				// 并返回预定义的错误 errno.AuthAccessExpired，表示访问令牌已过期
				return -1, errno.AuthAccessExpired
			}
			// 如果不是访问令牌，那么认为是刷新令牌过期
			// 使用 errno.NewErrNo 函数创建一个新的错误对象
			// 传入错误码 errno.AuthRefreshExpiredCode 和错误信息 "refresh token expired"
			// 同时返回 -1 作为处理后的令牌类型
			return -1, errno.NewErrNo(errno.AuthRefreshExpiredCode, "refresh token expired")
		}
	}

	// 如果错误不是 jwt.ValidationError 类型，或者不是令牌过期错误
	// 使用 fmt.Errorf 函数创建一个新的错误对象，包装原始错误 err
	// 错误信息为 "token validation failed"，并通过 %w 保留原始错误的详细信息
	// 同时返回 -1 作为处理后的令牌类型
	return -1, fmt.Errorf("token validation failed: %w", err)
}

// getTokenTTL 根据 token 类型返回过期时间
// 参数 tokenType 为令牌类型
// 返回值为过期时长
func getTokenTTL(tokenType int64) time.Duration {
	switch tokenType {
	case constants.TypeAccessToken:
		return constants.AccessTokenTTL
	case constants.TypeRefreshToken:
		return constants.RefreshTokenTTL
	default:
		return 0
	}
}
