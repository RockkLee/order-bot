package authsvc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"order-bot-mgmt-svc/internal/models"
)

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func signJWT(secret []byte, claims models.Claims) (string, error) {
	header := jwtHeader{Alg: "HS256", Typ: "JWT"}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("authsvc.signJWT: %w", err)
	}
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("authsvc.signJWT: %w", err)
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
		return models.Claims{}, fmt.Errorf("authsvc.parseJWT(), len(parts) != 3: %w", ErrInvalidToken)
	}
	tokenHeader, tokenPayload, tokenSignature := parts[0], parts[1], parts[2]
	enc := base64.RawURLEncoding
	signingInput := tokenHeader + "." + tokenPayload
	sign, err := enc.DecodeString(tokenSignature)
	if err != nil {
		return models.Claims{}, fmt.Errorf("authsvc.parseJWT(), failed to decode jwt signature: %w", ErrInvalidToken)
	}
	expectedSign := hmacSHA256(signingInput, secret)
	if !hmac.Equal(sign, expectedSign) {
		return models.Claims{}, fmt.Errorf("authsvc.parseJWT(), signature check failed : %w", ErrInvalidToken)
	}
	payloadBytes, err := enc.DecodeString(parts[1])
	if err != nil {
		return models.Claims{}, fmt.Errorf("authsvc.parseJWT(), failed to decode decoded header: %w", ErrInvalidToken)
	}
	var claims models.Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return models.Claims{}, fmt.Errorf("authsvc.parseJWT: %w", ErrInvalidToken)
	}
	if claims.Exp <= time.Now().Unix() {
		return models.Claims{}, fmt.Errorf("authsvc.parseJWT: %w", ErrExpiredToken)
	}
	return claims, nil
}

func hmacSHA256(message string, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(message))
	return mac.Sum(nil)
}
