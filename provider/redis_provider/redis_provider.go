package redis_provider

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"parking-service/provider/viper_provider"
	"time"
)

type IRedisProvider interface {
	SetKey(ctx context.Context , key string ,
		value interface{} , duration time.Duration) error
	GetKey(ctx context.Context , key string) (val string , err error)
	DeleteKey(ctx context.Context , key string)
}

type redisProvider struct {
	redisdb *redis.Client
}

func NewRedisProvider(IConfigProvider viper_provider.IConfigProvider) IRedisProvider {
	rdb := redis.NewClient(&redis.Options{
		Addr:     IConfigProvider.GetConfigEnv().RedisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &redisProvider{redisdb: rdb}
}

func(redis *redisProvider) SetKey(ctx context.Context , key string , 
	value interface{} , duration time.Duration) error {
	e , err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = redis.redisdb.Set(ctx, key, string(e), duration).Err()
	if err != nil {
		return err
	}

	return nil
}

func(redis *redisProvider) GetKey(ctx context.Context , key string) (val string , err error) {
	val, err = redis.redisdb.Get(ctx, key).Result()
	if err != nil {
		return
	}
	return
}

func(redis *redisProvider) DeleteKey(ctx context.Context , key string) {
	redis.redisdb.Del(ctx, key)
	return
}