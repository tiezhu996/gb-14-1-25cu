package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Claims JWT 载荷，包含用户 ID、角色、状态。
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT，有效期 24 小时。
func GenerateToken(secret string, userID primitive.ObjectID, username, role, status string) (string, error) {
	claims := Claims{
		UserID:   userID.Hex(),
		Username: username,
		Role:     role,
		Status:   status,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "codelearn",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并校验 JWT。
func ParseToken(secret, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
