package jwt

import (
	"errors"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

var (
	ErrMissingJWT       = errors.New("missing header authorization")
	ErrInvalidFormatJWT = errors.New("invalid header authorization format")
	ErrExpiredJWT       = errors.New("session expired")
	ErrInvalidGrade     = errors.New("invalid grade")
)

//nolint:tagliatelle
type SessionInfo struct {
	jwt.Claims

	Grade  string `json:"JWT_GRADE"`
	ID     string `json:"JWT_ID"`
	signer jose.Signer
}

func (s *SessionInfo) Signed() (string, error) {
	return jwt.Signed(s.signer).Claims(s).Serialize()
}

func (s *SessionInfo) CheckExpiration() error {
	if s.Expiry != nil && s.Expiry.Time().Before(time.Now()) {
		return ErrExpiredJWT
	}

	return nil
}

func (s *SessionInfo) Validate() error {
	if s.ID == "" {
		return ErrInvalidFormatJWT
	}

	return nil
}
