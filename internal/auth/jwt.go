package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("недопустимый токен")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func secret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))

}

func ttl() time.Duration {
	v := os.Getenv("JWT_TTL_MINUTES")
	if v == "" {
		return 60 * time.Minute
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return 60 * time.Minute
	}
	return time.Duration(n) * time.Minute
}

func GenerateToken(userID uint) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl())),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret())
}

func ParseToken(tokenStr string) (*Claims, error) {
	tkn, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return secret(), nil
	})
	if err != nil || !tkn.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := tkn.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
