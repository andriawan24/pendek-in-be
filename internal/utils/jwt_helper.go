package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserId    uuid.UUID `json:"user_id"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"

	defaultAccessTokenTTL  = 24 * time.Hour
	defaultRefreshTokenTTL = 30 * 24 * time.Hour

	defaultIssuer = "Link Short"
)

func tokenSecretFromEnv(key string) ([]byte, error) {
	secret := os.Getenv(key)
	if secret == "" {
		return nil, errors.New(key + " is not set")
	}
	return []byte(secret), nil
}

func generateToken(user database.User, tokenType string, ttl time.Duration, secretKey []byte) (string, JwtClaims, error) {
	claims := JwtClaims{
		UserId:    user.ID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    defaultIssuer,
			Subject:   user.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)

	return signedToken, claims, err
}

func GenerateJwtToken(user database.User) (string, JwtClaims, error) {
	return GenerateAccessToken(user)
}

func GenerateAccessToken(user database.User) (string, JwtClaims, error) {
	secret, err := tokenSecretFromEnv("TOKEN_SECRET")
	if err != nil {
		return "", JwtClaims{}, err
	}
	return generateToken(user, accessTokenType, defaultAccessTokenTTL, secret)
}

func GenerateRefreshToken(user database.User) (string, JwtClaims, error) {
	secret, err := tokenSecretFromEnv("REFRESH_TOKEN_SECRET")
	if err != nil {
		return "", JwtClaims{}, err
	}
	return generateToken(user, refreshTokenType, defaultRefreshTokenTTL, secret)
}

func ParseToken(tokenStr string) (*JwtClaims, error) {
	secret, err := tokenSecretFromEnv("TOKEN_SECRET")
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return secret, nil
	})

	if err != nil || token == nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token claims")
	}

	if claims.TokenType != "" && claims.TokenType != accessTokenType {
		return nil, errors.New("invalid token type")
	}

	return claims, err
}

func ParseRefreshToken(tokenStr string) (*JwtClaims, error) {
	secret, err := tokenSecretFromEnv("REFRESH_TOKEN_SECRET")
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || token == nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token claims")
	}

	if claims.TokenType != refreshTokenType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
