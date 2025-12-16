package helper

import (
	"auth/internal/model"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ParseExpiry(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty expiry string")
	}

	if strings.HasSuffix(s, "d") {
		nStr := strings.TrimSuffix(s, "d")
		n, err := strconv.Atoi(nStr)
		if err != nil {
			return 0, err
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return d, nil
}

func CreateSessionToken(user model.User) (string, error) {
	secretStr := os.Getenv("JWT_SECRET")
	expiryStr := os.Getenv("JWT_EXPIRED")

	if secretStr == "" {
		return "", errors.New("JWT_SECRET is not set in environment")
	}
	if expiryStr == "" {
		expiryStr = "30m"
	}

	duration, err := ParseExpiry(expiryStr)
	if err != nil {
		return "", err
	}

	expireAt := time.Now().Add(duration)

	claims := model.ClaimsModel{
		UserId:   user.ID,
		Role:     string(user.Role),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expireAt),
			Subject:   strconv.Itoa(user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretStr))
}

func CreateRefreshToken(user model.User) (string, error) {
	secretStr := os.Getenv("JWT_REFRESH_SECRET")
	expiryStr := os.Getenv("JWT_REFRESH_EXPIRED")

	if secretStr == "" {
		return "", errors.New("JWT_REFRESH_SECRET is not set in environment")
	}
	if expiryStr == "" {
		expiryStr = "7d"
	}

	duration, err := ParseExpiry(expiryStr)
	if err != nil {
		return "", err
	}

	expireAt := time.Now().Add(duration)

	claims := model.ClaimsModel{
		UserId:   user.ID,
		Role:     string(user.Role),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expireAt),
			Subject:   strconv.Itoa(user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretStr))
}

func VerifyToken(tokenString string) (*model.ClaimsModel, error) {
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		return nil, errors.New("JWT_SECRET missing")
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.ClaimsModel{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretStr), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.ClaimsModel); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func VerifyAndCreateNewSession(refreshTokenString string) (string, error) {
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		return "", errors.New("JWT_REFRESH_SECRET missing")
	}

	token, err := jwt.ParseWithClaims(refreshTokenString, &model.ClaimsModel{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to verify refresh token: %w", err)
	}

	claims, ok := token.Claims.(*model.ClaimsModel)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	user := model.User{
		ID:       claims.UserId,
		Username: claims.Username,
		Role:     model.Role(claims.Role),
	}

	newAccessToken, err := CreateSessionToken(user)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
