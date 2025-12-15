package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	// DefaultServiceName 默认的服务名
	DefaultServiceName = "go-sea"
	// DefaultTokenExpiration 默认的token过期时间
	DefaultTokenExpiration = 24 * time.Hour
	// DefaultSecretKey 默认的密钥
	DefaultSecretKey = "go-sea-internal-9k#L@m$pQ&x7*vB2nR5tY8uI3oP6wE1zA4sD"
)

// Claims 定义了JWT的声明
type Claims struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	jwt.RegisteredClaims
}

// GenerateServiceToken 生成一个用于服务间通信或内部操作的JWT
// userID: 用户标识
// serviceName: 服务名称，用于标识token的签发者
// expiration: token的有效期
// secretKey: 用于签名的密钥
func GenerateServiceToken(userID string, serviceName string, expiration time.Duration, secretKey string) (string, error) {
	if secretKey == "" {
		return "", fmt.Errorf("secret key cannot be empty")
	}

	// 创建声明
	claims := Claims{
		UserID:      userID,
		ServiceName: serviceName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    serviceName,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	// 使用HS256签名算法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的secret签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateServiceTokenWithDefaultExpiry 使用默认的过期时间生成token
func GenerateServiceTokenWithDefaultExpiry(userID string, serviceName string, secretKey string) (string, error) {
	return GenerateServiceToken(userID, serviceName, DefaultTokenExpiration, secretKey)
}

// GenerateServiceTokenWithDefaults 使用默认的服务名和过期时间生成token
func GenerateServiceTokenWithDefaults(userID string, secretKey string) (string, error) {
	return GenerateServiceToken(userID, DefaultServiceName, DefaultTokenExpiration, secretKey)
}

// ParseServiceToken 解析服务间通信的JWT
func ParseServiceToken(tokenString string, secretKey string) (*Claims, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("secret key cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保token方法是预期的加密方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

// GenerateServiceTokenDefault 使用默认密钥生成服务token
func GenerateServiceTokenDefault(userID string, serviceName string, expiration time.Duration) (string, error) {
	return GenerateServiceToken(userID, serviceName, expiration, DefaultSecretKey)
}

// GenerateServiceTokenDefaultWithDefaults 使用默认密钥、默认服务名和默认过期时间生成token
func GenerateServiceTokenDefaultWithDefaults(userID string) (string, error) {
	return GenerateServiceToken(userID, DefaultServiceName, DefaultTokenExpiration, DefaultSecretKey)
}

// ParseServiceTokenDefault 使用默认密钥解析服务token
func ParseServiceTokenDefault(tokenString string) (*Claims, error) {
	return ParseServiceToken(tokenString, DefaultSecretKey)
}
