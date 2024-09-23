package jwt

import (
	"auth/internal/models"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	genTime    time.Time
)

func MustSetKeys(log *slog.Logger) {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Error("failed to generate private key", slog.Any("error", err))
		os.Exit(1)
	}

	publicKey = &privateKey.PublicKey

	genTime = time.Now()

	log.Info("keys are generated", slog.Time("generation time", genTime))
}

func NewToken(usr *models.User, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"login":    usr.Login,
		"password": password,
		"id":       usr.Id,
	})

	return token.SignedString(privateKey)
}

func GetPublicKey() []byte {
	b := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}

	return pem.EncodeToMemory(b)
}

func GetKeysGenerationTime() time.Time {
	return genTime
}
