package storage

import (
	"auth/internal/models"
	"auth/internal/storage/postgres"
	"auth/internal/storage/redis"
	"context"
	"errors"
	"fmt"
)

type Storage struct {
	db    postgres.Database
	cache redis.Cache
}

var (
	ErrRedis = errors.New("redis error")
)

func MustNew(dbConnStr string, redisAddr string) Storage {
	return Storage{
		db:    postgres.New(dbConnStr),
		cache: redis.New(redisAddr),
	}
}

func (s Storage) SaveUser(ctx context.Context, login string, passHash []byte) (*models.User, error) {
	u, err := s.db.InsertUser(login, passHash)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SaveUser(u); err != nil {
		return u, fmt.Errorf("%w: %w", ErrRedis, err)
	}
	return u, nil
}

func (s Storage) GetUser(ctx context.Context, login string) (*models.User, error) {
	u, e := s.cache.GetUser(login)
	if e == nil {
		return u, nil
	}

	if errors.Is(e, redis.ErrNotFound) {
		e = nil
	} else {
		e = fmt.Errorf("%w: %w", ErrRedis, e)
	}

	u, err := s.db.GetUser(login)
	if err != nil {
		return nil, err
	}

	return u, e
}

func (s Storage) Close() error {
	e := s.cache.Close()
	if e != nil {
		err := s.db.Close()
		return fmt.Errorf("%w, %w", e, err)
	} else {
		return s.db.Close()
	}
}
