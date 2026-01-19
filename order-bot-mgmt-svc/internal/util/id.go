package util

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func NewID() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(buf)
}
