package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

func GenerateSecretKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		logger.Log.Error("Error to generate secret key", zap.Error(err))
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func GenerateTokenValue(userId int64, key string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userId,
		"expiration": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		logger.Log.Error("Error to generate jwt", zap.Error(err))
		return "", err
	}

	return signedToken, nil
}

func ValidateTokenAndGetUserId(tokenValue string, key string) (int64, error) {
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return -1, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return -1, errors.New("incorrect token")
	}

	userIdFloat, ok := claims["user_id"].(float64)
	if !ok {
		return -1, errors.New("incorrect user_id field inside token")
	}
	userId := int64(userIdFloat)

	if exp, ok := claims["expiration"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		if !expTime.After(time.Now()) {
			logger.Log.Info("Token expired", zap.Int64("user_id", userId))
			return -1, errors.New("token expired")
		}
	} else {
		return -1, errors.New("error to get token expiration time")
	}

	return userId, nil
}
