package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	prefix string
}

func NewRedisCache(host string, port int, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("RedisCache: 连接Redis失败(%s:%d): %v，将跳过缓存", host, port, err)
		return &RedisCache{}, nil
	}

	log.Printf("RedisCache: 连接Redis成功 (%s:%d)", host, port)
	return &RedisCache{client: client, prefix: "ra:"}, nil
}

func (c *RedisCache) IsAvailable() bool {
	return c.client != nil
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	if c.client == nil {
		return false, nil
	}

	val, err := c.client.Get(ctx, c.prefix+key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}

	return true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.prefix+key, data, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return nil
	}

	return c.client.Del(ctx, c.prefix+key).Err()
}

func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	if c.client == nil {
		return nil
	}

	keys, err := c.client.Keys(ctx, c.prefix+pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

func (c *RedisCache) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
