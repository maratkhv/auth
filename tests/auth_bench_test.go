package tests

import (
	"auth/internal/lib/jwt"
	"auth/internal/models"
	"auth/internal/services/auth"
	m "auth/tests/mocks"
	"bytes"
	"context"
	"log/slog"
	"testing"

	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func setupApi(b *testing.B, ctx context.Context, login string, password string) *auth.Auth {
	ctrl := gomock.NewController(b)

	up := m.NewMockUserProvider(ctrl)
	us := m.NewMockUserSaver(ctrl)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(
		&logBuffer,
		&slog.HandlerOptions{},
	))

	jwt.MustSetKeys(logger)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), auth.HashCost)
	if err != nil {
		b.Fatalf("error generating password hash: %v", err)
	}

	mockUser := &models.User{
		Id:       1,
		Login:    login,
		PassHash: passHash,
	}

	up.EXPECT().GetUser(gomock.All(), gomock.All()).Return(mockUser, nil).AnyTimes()
	us.EXPECT().SaveUser(gomock.All(), gomock.All(), gomock.All()).Return(mockUser, nil).AnyTimes()

	return auth.New(logger, up, us)
}

func BenchmarkRegister(b *testing.B) {
	ctx := context.Background()
	login := "some-login"
	pwd := "some-password"

	api := setupApi(b, ctx, login, pwd)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		api.Register(ctx, login, pwd)
	}
}

func BenchmarkLogin(b *testing.B) {
	ctx := context.Background()
	login := "some-login"
	pwd := "some-password"

	api := setupApi(b, ctx, login, pwd)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		api.Login(ctx, login, pwd)
	}
}
