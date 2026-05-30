// Package jwt 提供 JWT（JSON Web Token）的生成与解析功能。
// JWT 由三部分组成：
//   - Header（头部）：声明令牌类型（typ: "JWT"）和签名算法（alg: "HS256"）
//   - Payload（载荷）：存放自定义声明（如 userId、username、roles）和标准声明（如过期时间）
//   - Signature（签名）：使用密钥对 Header + Payload 进行签名，防止数据被篡改
// 本包使用 HS256（HMAC-SHA256）对称签名算法，即同一个密钥既用于签名也用于验证。
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 预定义错误变量（sentinel errors），用于区分 JWT 验证失败的不同原因。
// Go 语言中通过 errors.New() 创建哨兵错误，上层调用方可以使用 errors.Is() 判断具体错误类型。
var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("无效的token")
	ErrTokenMalformed   = errors.New("token格式错误")
	ErrTokenNotValidYet = errors.New("token尚未生效")
)

// Claims 定义 JWT Payload 中的自定义声明（Custom Claims）字段。
// 内嵌 jwt.RegisteredClaims 提供了标准声明字段：
//   - ExpiresAt（过期时间）：Token 超过此时间后失效
//   - IssuedAt（签发时间）：Token 的签发时间
//   - NotBefore（生效时间）：Token 在此时间之前不可用
//   - Issuer（签发者）：标识 Token 的签发方
// 这种结构体嵌套（embedding）是 Go 的惯用方式，子结构体的字段可直接通过外层访问。
type Claims struct {
	UserId      uint64   `json:"user_id"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// RefreshClaims holds the refresh token payload.
type RefreshClaims struct {
	jwt.RegisteredClaims
	UserId uint64 `json:"user_id"`
}

// JWT 持有签名密钥、过期时间和签发者信息，是生成与解析 Token 的核心结构体。
// 字段使用小写字母开头，表示包内私有（unexported），外部只能通过 New() 函数创建。
type JWT struct {
	secret    []byte         // 签名密钥，HS256 使用同一密钥签名和验证
	expires   time.Duration  // Token 有效期，如 24 小时
	issuer    string         // 签发者标识，如 "tool-go"
}

// New 创建 JWT 实例的构造函数（constructor）。
// Go 中没有类（class），通过 New 开头的函数模拟构造函数，返回结构体指针 *JWT。
// secret 为签名密钥，expires 为 Token 有效期，issuer 为签发者标识。
func New(secret string, expires time.Duration, issuer string) *JWT {
	return &JWT{
		secret: []byte(secret),
		expires: expires,
		issuer: issuer,
	}
}

// GenerateToken 根据用户信息生成 JWT Token 字符串。
// 使用 HS256（HMAC-SHA256）算法进行签名。HS256 是对称签名算法：
//   - 签名：用密钥对 Header + Payload 计算 HMAC-SHA256，附加为 Signature
//   - 验证：接收方用相同的密钥重新计算签名，比对是否一致
// 对称算法的优点是计算速度快，缺点是需要安全地共享密钥。
func (j *JWT) GenerateToken(userId uint64, username string, roles []string, permissions []string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:      userId,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expires)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ParseToken 解析并验证 JWT Token，返回 Claims。
// 验证过程包括：
//   1. 格式校验：Token 是否为三段式结构（Header.Payload.Signature）
//   2. 签名验证：使用密钥重新计算签名，比对是否匹配
//   3. 时间验证：检查是否过期（ExpiresAt）、是否已生效（NotBefore）
// 验证失败时根据具体原因返回对应的预定义错误，便于调用方做差异化处理。
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotValidYet
		default:
			return nil, ErrTokenInvalid
		}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// GenerateRefreshToken creates a long-lived refresh token (7 days).
func (j *JWT) GenerateRefreshToken(userId uint64) (string, error) {
	claims := RefreshClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expires * 28)), // ~7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ParseRefreshToken validates and parses a refresh token string.
func (j *JWT) ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			}
		}
		return nil, ErrTokenInvalid
	}
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}
