package jwtutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"order-bot-mgmt-svc/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func SignJWT(secret []byte, claims models.Claims) (string, error) {
	header := jwtHeader{Alg: "HS256", Typ: "JWT"}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("jwtutil.SignJWT: %w", err)
	}
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("jwtutil.SignJWT: %w", err)
	}
	enc := base64.RawURLEncoding
	headerB64 := enc.EncodeToString(headerBytes)
	payloadB64 := enc.EncodeToString(payloadBytes)
	signingInput := headerB64 + "." + payloadB64
	signature := hmacSHA256(signingInput, secret)
	sigB64 := enc.EncodeToString(signature)
	return signingInput + "." + sigB64, nil
}

func ParseJWT(secret []byte, token string) (models.Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return models.Claims{}, fmt.Errorf("jwtutil.ParseJWT(), len(parts) != 3: %w", ErrInvalidToken)
	}
	tokenHeader, tokenPayload, tokenSignature := parts[0], parts[1], parts[2]
	enc := base64.RawURLEncoding
	signingInput := tokenHeader + "." + tokenPayload
	sign, err := enc.DecodeString(tokenSignature)
	if err != nil {
		return models.Claims{}, fmt.Errorf("jwtutil.ParseJWT(), failed to decode jwt signature: %w", ErrInvalidToken)
	}
	expectedSign := hmacSHA256(signingInput, secret)
	if !hmac.Equal(sign, expectedSign) {
		return models.Claims{}, fmt.Errorf("jwtutil.ParseJWT(), signature check failed : %w", ErrInvalidToken)
	}
	payloadBytes, err := enc.DecodeString(parts[1])
	if err != nil {
		return models.Claims{}, fmt.Errorf("jwtutil.ParseJWT(), failed to decode decoded header: %w", ErrInvalidToken)
	}
	var claims models.Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return models.Claims{}, fmt.Errorf("jwtutil.ParseJWT: %w", ErrInvalidToken)
	}
	if claims.Exp <= time.Now().Unix() {
		return models.Claims{}, fmt.Errorf("jwtutil.ParseJWT: %w", ErrExpiredToken)
	}
	return claims, nil
}

func hmacSHA256(message string, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(message))
	return mac.Sum(nil)
}

func GetToken(w http.ResponseWriter, r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return "", true
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return "", true
	}
	accessToken := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	return accessToken, false
}

func GetTokenGin(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", false
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", false
	}
	accessToken := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	return accessToken, accessToken != ""
}
