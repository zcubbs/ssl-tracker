package token

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

const minSecretKeySize = 32

type JwtMaker struct {
	secretKey string
}

// NewJwtMaker creates a new JWT maker
func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrInvalidSecretKey
	}
	return &JwtMaker{
		secretKey: secretKey,
	}, nil
}

// CreateToken creates a new token for a specific username and duration.
func (maker *JwtMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

// VerifyToken checks if the token is valid or not.
func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok { // if not ok, return error
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		var verr *jwt.ValidationError
		ok := errors.As(err, &verr)
		if ok && (verr.Errors&jwt.ValidationErrorExpired != 0) { // if token is expired
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
