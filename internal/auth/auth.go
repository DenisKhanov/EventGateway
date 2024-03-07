// Package auth provides functions for handling authentication, JWT token creation,
// and validation.
package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"

	"github.com/sirupsen/logrus"
	"time"
)

// Claims is a structure that includes standard JWT claims and UserID.
type Claims struct {
	jwt.RegisteredClaims
	TopicName string
}

// const for generate token
const (
	// TokenExp defines the expiration duration for JWT tokens.
	TokenExp = time.Hour * 24
	// SecretKey is the secret key used for signing JWT tokens.
	SecretKey = "SnJSkf123jlLKNfsNln"
)

// BuildJWTString creates a token with the HS256 signature algorithm and Claims statements and returns it as a string.
func BuildJWTString(topicName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		TopicName: topicName,
	})
	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return tokenString, nil
}

// IsValidToken method to check the token for validity, we return bool
func IsValidToken(tokenString string) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		logrus.Error(err)
		return false
	}
	if !token.Valid {
		err = fmt.Errorf("token is not valid")
		logrus.Error(err)
		return false
	}
	return true
}

// GetUserTopicName we check the validity of the token and if it is valid, then we get and return the UserID from it
func GetUserTopicName(tokenString string) (string, error) {
	fmt.Println(tokenString)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	if !token.Valid {
		err = fmt.Errorf("token is not valid")
		logrus.Error(err)
		return "", err
	}
	logrus.Infof("Token is valid, topic name: %v", claims.TopicName)
	return claims.TopicName, nil
}
