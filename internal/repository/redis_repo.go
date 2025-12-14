package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

)

// RedisRepository diubah menjadi huruf besar (Public)
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository sekarang mengembalikan *RedisRepository (Struct Pointer)
// Binding ke interface domain.RedisRepository dilakukan di wire.go
func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisRepository) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisRepository) DeletePrefix(ctx context.Context, prefix string) error {
	// Menggunakan SCAN (bukan KEYS) agar aman untuk production database yang besar
	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()

	for iter.Next(ctx) {
		// Hapus satu per satu key yang ditemukan
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}