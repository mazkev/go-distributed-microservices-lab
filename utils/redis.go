package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client    *redis.Client
	isOnline  bool
}

var Redis *RedisClient

// InitRedis menginisialisasi koneksi ke server Redis
func InitRedis(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	online := (err == nil)

	if online {
		fmt.Println("✅ Terhubung ke Redis Server di", addr)
	} else {
		fmt.Printf("⚠️ Redis belum aktif di %s (API tetap berjalan dengan fallback langsung ke database)\n", addr)
	}

	Redis = &RedisClient{
		client:   rdb,
		isOnline: online,
	}
	return Redis
}

// Get mengambil dan meng-unmarshal data dari cache Redis
func (r *RedisClient) Get(ctx context.Context, key string, target any) bool {
	if r == nil || !r.isOnline {
		return false
	}

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(val), target)
	return err == nil
}

// Set menyimpan data ke Redis dalam format JSON dengan durasi TTL
func (r *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if r == nil || !r.isOnline {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// Del menghapus satu atau beberapa key di Redis (Cache Invalidation)
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	if r == nil || !r.isOnline {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}
