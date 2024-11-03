package jwt

import (
	"auth/internal/models"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestKeysGeneration(t *testing.T) {
	setKeys()
	key := GetPublicKey()

	if key == nil || len(key) == 0 {
		t.Fatal("failed to get key")
	}

	b, _ := pem.Decode(key)
	if b == nil {
		t.Fatal("failed to decode key")
	}

	_, err := x509.ParsePKCS1PublicKey(b.Bytes)
	if err != nil {
		t.Fatalf("failed to parse key: %v", err)
	}
}

func TestToken(t *testing.T) {
	setKeys()

	tests := []struct {
		name string
		usr  *models.User
		pwd  string
	}{
		{"normal", &models.User{Id: 1, Login: "user1"}, "123"},
		{"normal2", &models.User{Id: 12, Login: "user12"}, "1233"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tokenStr, err := NewToken(tc.usr, tc.pwd)
			if err != nil {
				t.Fatalf("failed to get token: %v", err)
			}

			key, err := jwt.ParseRSAPublicKeyFromPEM(GetPublicKey())
			if err != nil {
				t.Fatalf("failed to parse key: %v", err)
			}

			token, err := jwt.Parse(tokenStr, func(tk *jwt.Token) (interface{}, error) {
				if _, ok := tk.Method.(*jwt.SigningMethodRSA); !ok {
					t.Fatalf("signing method is not rsa: %T", tk.Method)
				}
				return key, nil
			})
			if err != nil {
				t.Fatalf("failed to parse token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				t.Fatalf("wrong claims type: %T instead of jwt.MapClaims", token.Claims)
			}

			assert.Equal(t, tc.usr.Id, int64(claims["id"].(float64)))
			assert.Equal(t, tc.pwd, claims["password"])
			assert.Equal(t, tc.usr.Login, claims["login"])
		})
	}
}

func setKeys() bytes.Buffer {
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(
		&logBuffer,
		&slog.HandlerOptions{},
	))

	MustSetKeys(logger)
	return logBuffer
}
