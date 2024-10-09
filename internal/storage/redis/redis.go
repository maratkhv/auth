package redis

import (
	"auth/internal/models"
	"context"
	"encoding/json"
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
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.c.Set(context.Background(), "user@"+user.Login, bytes, time.Hour).Err()
}

func (c *Cache) GetUser(login string) (*models.User, error) {
	bytes := make([]byte, 0)
	err := c.c.Get(context.Background(), "user@"+login).Scan(&bytes)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var u models.User
	err = json.Unmarshal(bytes, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (c *Cache) Close() error {
	return c.c.Close()
}
