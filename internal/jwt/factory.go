package jwt

import (
	"encoding/hex"
	"time"

	"irj/pkg/utils"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

func NewSigner(secret string) jose.Signer {
	key, err := hex.DecodeString(secret)
	if err != nil {
		return nil
	}

	options := jose.SignerOptions{}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, options.WithType("JWT"))
	if err != nil {
		return nil
	}

	return signer
}

func NewToken(signer jose.Signer, sessionDuration time.Duration, grade, id string) SessionInfo {
	iat := time.Now()

	tokID, _ := utils.CreateUUID()

	return SessionInfo{
		Claims: jwt.Claims{
			Issuer:    "IRJ",
			Expiry:    jwt.NewNumericDate(iat.Add(sessionDuration)),
			NotBefore: jwt.NewNumericDate(iat),
			IssuedAt:  jwt.NewNumericDate(iat),
			ID:        tokID,
		},
		ID:     id,
		Grade:  grade,
		signer: signer,
	}
}

func CheckSignatureAndDecode(rawToken string, key []byte) (SessionInfo, error) {
	if rawToken == "" {
		return SessionInfo{}, ErrMissingJWT
	}

	tok, err := jwt.ParseSigned(rawToken, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return SessionInfo{}, ErrInvalidFormatJWT
	}

	claims := SessionInfo{}

	jwk := jose.JSONWebKey{
		Key: key,
	}

	if err = tok.Claims(jwk, &claims); err != nil {
		return SessionInfo{}, ErrInvalidFormatJWT
	}

	return claims, nil
}
