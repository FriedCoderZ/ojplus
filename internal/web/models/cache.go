package models

import (
	"Alarm/internal/config"
	"errors"

	"github.com/go-redis/redis"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(cfg *config.Redis) (*Cache, error) {
	if cfg.Addr == "" {
		return nil, errors.New("the address must be given")
	}
	cache := Cache{}
	cache.Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &cache, nil
}
