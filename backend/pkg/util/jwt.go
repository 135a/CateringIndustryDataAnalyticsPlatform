package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 实际项目中建议将 Secret 存入 .env 环境变量中
var jwtSecret = []byte("catering_secret_key_2026")

// Claims 包含我们想要封装在 Token 中的载荷数据
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token，有效期设置为 24 小时
func GenerateToken(userID int64, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)

	claims := Claims{
		userID,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "catering-platform",
		},
	}

	// 使用 HS256 算法进行签名
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}
