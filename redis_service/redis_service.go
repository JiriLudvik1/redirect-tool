package redis_service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisService(addr string, password string, db int) (*RedisService, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis service connection failed: %s", err)
		return nil, err
	}

	return &RedisService{
		client: rdb,
		ctx:    ctx,
	}, nil
}

func (s *RedisService) CreateRedirectEntry(url string) (string, error) {
	urlHash := createUrlHash(url)
	redirectKey := getRedirectRedisKey(urlHash)

	err := s.set(redirectKey, url)
	if err != nil {
		return "", nil
	}

	return urlHash, nil
}

func (s *RedisService) GetOriginalUril(hash string) (string, error) {
	redirectKey := getRedirectRedisKey(hash)

	value, err := s.get(redirectKey)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *RedisService) set(key string, values string) error {
	return s.client.Set(s.ctx, key, values, 0).Err()
}

func (s *RedisService) get(key string) (string, error) {
	return s.client.Get(s.ctx, key).Result()
}
