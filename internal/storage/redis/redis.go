package redis

import (
	"auth/internal/models"
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	c *redis.Client
}

var (
	ErrNotFound = errors.New("user not found")
)

func New(addr string) Cache {
	return Cache{c: redis.NewClient(&redis.Options{
		Addr:       addr,
		MaxRetries: -1,
	})}
}

func (c *Cache) SaveUser(user *models.User) error {
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(*user); err != nil {
		return err
	}

	return c.c.Set(context.Background(), "user@"+user.Login, buf.Bytes(), time.Hour).Err()
}

func (c *Cache) GetUser(login string) (*models.User, error) {
	bs := make([]byte, 0)
	err := c.c.Get(context.Background(), "user@"+login).Scan(&bs)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var u models.User
	buf := bytes.NewBuffer(bs)
	if err := gob.NewDecoder(buf).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (c *Cache) Close() error {
	return c.c.Close()
}
