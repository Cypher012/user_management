package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func HashTokenKey(rawkey, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(rawkey))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRandomString(length int) (string, error) {
	// base64 expands by 4/3, so we reverse it
	byteLen := (length * 3) / 4
	b := make([]byte, byteLen)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

func GenerateToken() (string, error) {
	s, err := GenerateRandomString(32)
	if err != nil {
		return "", err
	}
	return "etk" + s, nil
}
