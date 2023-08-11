package token

import (
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new Paseto maker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < minSecretKeySize {
		return nil, ErrInvalidSecretKey
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

// CreateToken creates a new token for a specific username and duration.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

// VerifyToken checks if the token is valid or not.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
