package jwt_test

import (
	"auth/internal/lib/jwt"
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"testing"
)

func TestGetPublicKey(t *testing.T) {
	jwt.MustSetKeys(slog.Default())
	key := jwt.GetPublicKey()

	if key == nil || len(key) == 0 {
		t.Error("failed to get key")
	}

	b, _ := pem.Decode(key)
	if b == nil {
		t.Error("failed to decode key")
	}

	_, err := x509.ParsePKCS1PublicKey(b.Bytes)
	if err != nil {
		t.Error(err)
	}

}
