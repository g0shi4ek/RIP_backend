package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	blacklistPrefix    = "blacklist:%s"
)

type RedisClient struct {
	Client *redis.Client
}

type RedisConfig struct {
	Password    string
	Username    string
	Endpoint    string
	DB          int
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

func RedisConfigFromEnv() (*RedisConfig, error) {
	pass := os.Getenv("REDIS_PASSWORD")
	user := os.Getenv("REDIS_USER")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	dialTimeout, _ := time.ParseDuration(os.Getenv("REDIS_DIAL_TIMEOUT"))
	readTimeout, _ := time.ParseDuration(os.Getenv("REDIS_READ_TIMEOUT"))
	endpoint := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	return &RedisConfig{
		Password:    pass,
		Username:    user,
		Endpoint:    endpoint,
		DB:          db,
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
	}, nil
}

func NewRedisClient() (*RedisClient, error) {
	cfg, err := RedisConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load redis config")
	}
	client := redis.NewClient(&redis.Options{
		Addr:        cfg.Endpoint,
		Password:    cfg.Password,
		Username:    cfg.Username,
		DB:          cfg.DB,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &RedisClient{
		Client: client,
	}, nil
}

func (c *RedisClient) Close() error {
	return c.Client.Close()
}

func (c *RedisClient) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := fmt.Sprintf(blacklistPrefix, token)
	return c.Client.Set(ctx, key, true, ttl).Err()
}

func (c *RedisClient) IsInBlacklist(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf(blacklistPrefix, token)
	result, err := c.Client.Exists(ctx, key).Result()
	return result > 0, err
}