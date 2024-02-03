package token

import (
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	minSecretKeySize = 32
)

var (
	InvalidKeySize  = fmt.Errorf("ERROR: secret key size must greater than %d", minSecretKeySize)
	InvalidTokenAlg = fmt.Errorf("unexpected signing method, expected SigningMethodHS256")
	InvalidToken    = fmt.Errorf("invalid token")
)

// JWTMaker is a Json web token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker create new jwt maker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, InvalidKeySize
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken create a new jwt token
func (J JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	key, err := base64.StdEncoding.DecodeString(J.secretKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

// VerifyToken verify if the jwt token is valid or not
func (J JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, InvalidTokenAlg
		}

		key, err := base64.StdEncoding.DecodeString(J.secretKey)
		if err != nil {
			return "", err
		}

		return key, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, InvalidToken
	}

	return payload, nil
}
