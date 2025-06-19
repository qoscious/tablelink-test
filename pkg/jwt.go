package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	secretKey string
	expired   time.Duration
}

func NewJWTManager(secretKey string, expired time.Duration) *JWTManager {
	return &JWTManager{secretKey, expired}
}

type UserClaims struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
	jwt.RegisteredClaims
}

func (j *JWTManager) Generate(userID, roleID string) (string, error) {
	claims := UserClaims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expired)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) Verify(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
