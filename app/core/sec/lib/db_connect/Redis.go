package db_connect

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func CreateRedisConnection(dbconf RedisConfiguration) (*RedisClientWrapper, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     dbconf.Host + ":" + dbconf.Port,
		Password: dbconf.Password,
		DB:       dbconf.ID,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisClientWrapper{Client: client}, err
}

type RedisConfiguration struct {
	Host     string
	Port     string
	Password string
	ID       int
}



type RedisClientWrapper struct {
	*redis.Client
}

func (c *RedisClientWrapper) SetStruct(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	p, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return c.Client.Set(ctx, key, p, expiration).Result()
}

func (c *RedisClientWrapper) GetUnmarshal(ctx context.Context, key string, stk interface{}) error {
	p, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	//fmt.Println("p: ",p)
	if p != "" {
		e := json.Unmarshal([]byte(p), stk)
		if e != nil {
			return fmt.Errorf("error when json.Unmarshal: %w", e)
		}
		return nil
	} else {
		return nil
	}
}
