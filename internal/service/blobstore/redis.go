package blobstore

import (
	"bytes"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"math/rand"
	"time"
)

var _ TempBlobStore = (*Redis)(nil)

type Redis struct {
	rdb      *redis.Client
	lifetime time.Duration
}

type Config struct {
	FilesLifetime time.Duration
	RedisUrl      string
}

func InitRedis(cfg Config) (*Redis, error) {
	opts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, fmt.Errorf("parse redis url from cfg: %s", err)
	}
	rdb := redis.NewClient(opts)
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("ping redis client: %s", err)
	}
	return &Redis{
		rdb:      rdb,
		lifetime: cfg.FilesLifetime,
	}, nil
}

func (r *Redis) Save(data io.Reader, extension string) (*URI, error) {
	b, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("read from input data")
	}
	uri := generateURI() + "." + extension
	ctx := context.Background()
	err = r.rdb.Set(ctx, uri, b, r.lifetime).Err()
	if err != nil {
		return nil, fmt.Errorf("set value to redis: %s", err)
	}
	return &uri, nil
}

func (r *Redis) Get(uri URI) (io.Reader, error) {
	res := r.rdb.Get(context.Background(), uri)
	if res.Err() != nil {
		return nil, fmt.Errorf("get b from redis: %s", res.Err())
	}
	b, err := res.Bytes()
	if err != nil {
		return nil, fmt.Errorf("convert b to bytes: %s", res.Err())
	}
	rc := bytes.NewReader(b)
	return rc, err
}

// generateURI генерирует URI на основе текущей даты и случайной строки
func generateURI() URI {
	randomString := randString(5)                           // Генерируем случайную строку
	currentDate := time.Now().Format("2006-01-02-15-04-05") // Форматируем текущую дату
	return currentDate + "_" + randomString
}

// Генерация случайной строки
func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
