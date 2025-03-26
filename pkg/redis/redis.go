package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func New(addr, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Redis{Client: client}, nil
}

func (r *Redis) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}) error {
	data, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	} else if err != nil {
		return ErrCacheFailed
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return ErrUnmarshal
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrMarshal
	}

	if err := r.Client.Set(ctx, key, data, 0).Err(); err != nil {
		return ErrCacheFailed
	}

	return nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		return ErrCacheFailed
	}

	return nil
}

func (r *Redis) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return ErrCacheFailed
	}

	if len(keys) == 0 {
		return ErrCacheMiss
	}

	if err := r.Client.Del(ctx, keys...).Err(); err != nil {
		return ErrCacheFailed
	}

	return nil
}
