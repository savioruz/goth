package jwt

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	instance *JWT
	once     sync.Once
)

type JWT struct {
	appName            string
	secretKey          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func Initialize(appName string, secretKey string, accessExpiry, refreshExpiry time.Duration) {
	once.Do(func() {
		instance = &JWT{
			appName:            appName,
			secretKey:          secretKey,
			accessTokenExpiry:  accessExpiry,
			refreshTokenExpiry: refreshExpiry,
		}
	})
}

func GetInstance() *JWT {
	if instance == nil {
		panic(errors.New("jwt not initialized"))
	}
	return instance
}

func GenerateAccessToken(userID, email, level string) (string, error) {
	return GetInstance().generateToken(userID, email, level, GetInstance().accessTokenExpiry, "access_token")
}

func GenerateRefreshToken(userID, email, level string) (string, error) {
	return GetInstance().generateToken(userID, email, level, GetInstance().refreshTokenExpiry, "refresh_token")
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetInstance().secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWT) generateToken(userID, email, level string, expiry time.Duration, tokenType string) (string, error) {
	claims := &Claims{
		ID:        userID,
		Email:     email,
		Level:     level,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.appName,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.secretKey))
}
