package auth

import (
	"auth/internal/lib/jwt"
	"auth/internal/models"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrProvider UserProvider
	usrSaver    UserSaver
}

var (
	ErrWrongPassword = errors.New("wrong password")
)

type UserProvider interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type UserSaver interface {
	SaveUser(ctx context.Context, login string, passHash []byte) (*models.User, error)
}

func New(log *slog.Logger, prov UserProvider, usrSaver UserSaver) *Auth {
	return &Auth{
		log:         log,
		usrProvider: prov,
		usrSaver:    usrSaver,
	}
}

func (a *Auth) Login(ctx context.Context, login string, password string) (token string, keySinceUnix int64, err error) {
	usr, err := a.usrProvider.GetUser(ctx, login)
	if err != nil {
		if !errors.Is(err, storage.ErrRedis) {
			a.log.Error("error getting user", slog.Any("error", err))
			return "", 0, err
		} else {
			a.log.Warn("redis error", slog.Any("error", err))
		}
	}

	if err := bcrypt.CompareHashAndPassword(usr.PassHash, []byte(password)); err != nil {
		a.log.Warn("wrong password", slog.Any("error", err))
		return "", 0, ErrWrongPassword
	}

	token, err = jwt.NewToken(usr, password)
	if err != nil {
		return "", 0, err
	}

	a.log.Debug("created new token", slog.String("token", token))
	t := jwt.GetKeysGenerationTime().Unix()

	return token, t, nil
}

func (a *Auth) Register(ctx context.Context, login string, password string) (id int64, err error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("error generating hash from password", slog.Any("error", err))
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	usr, err := a.usrSaver.SaveUser(ctx, login, passHash)
	if err != nil {
		a.log.Error("failed to save user", slog.Any("error", err))
		return 0, fmt.Errorf("failed to register user: %w", err)
	}

	return usr.Id, nil
}

func (a *Auth) GetPublicKey(ctx context.Context) (key []byte, createTime time.Time, err error) {
	key = jwt.GetPublicKey()
	if key == nil {
		a.log.Error("failed to encode public key")
		return nil, createTime, errors.New("internal error")
	}
	createTime = jwt.GetKeysGenerationTime()
	return key, createTime, nil
}
