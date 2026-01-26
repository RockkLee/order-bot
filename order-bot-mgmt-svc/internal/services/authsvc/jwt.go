package authsvc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"order-bot-mgmt-svc/internal/models"
)

var errInvalidToken = errors.New("invalid token")

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func signJWT(secret []byte, claims models.Claims) (string, error) {
	header := jwtHeader{Alg: "HS256", Typ: "JWT"}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	enc := base64.RawURLEncoding
	headerB64 := enc.EncodeToString(headerBytes)
	payloadB64 := enc.EncodeToString(payloadBytes)
	signingInput := headerB64 + "." + payloadB64
	signature := hmacSHA256(signingInput, secret)
	sigB64 := enc.EncodeToString(signature)
	return signingInput + "." + sigB64, nil
}

func parseJWT(secret []byte, token string) (models.Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return models.Claims{}, errInvalidToken
	}
	enc := base64.RawURLEncoding
	signingInput := parts[0] + "." + parts[1]
	sig, err := enc.DecodeString(parts[2])
	if err != nil {
		return models.Claims{}, errInvalidToken
	}
	expectedSig := hmacSHA256(signingInput, secret)
	if !hmac.Equal(sig, expectedSig) {
		panic("expectedSig != sig")
		return models.Claims{}, errInvalidToken
	}
	payloadBytes, err := enc.DecodeString(parts[1])
	if err != nil {
		return models.Claims{}, errInvalidToken
	}
	var claims models.Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return models.Claims{}, errInvalidToken
	}
	if claims.Exp <= time.Now().Unix() {
		return models.Claims{}, errInvalidToken
	}
	return claims, nil
}

func hmacSHA256(message string, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(message))
	return mac.Sum(nil)
}
