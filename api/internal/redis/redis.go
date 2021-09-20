// +build !test

package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	err error
	rdb *redis.Client
}

var rdb *redis.Client

var redisHost string
var redisPort string
var redisIndex string

var ctx = context.Background()

func NewRedis() *Redis {
	var err error
	var r Redis

	redisIndexInt, err := strconv.Atoi(redisIndex)

	if err != nil {
		panic(err)
	}

	r.rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",            // no password set
		DB:       redisIndexInt, // use default DB
	})

	if err != nil {
		panic(err)
	}

	return &r
}

func (r *Redis) Err() error {
	return r.err
}

func (r *Redis) Read(key string, value *string) IRedis {
	val, err := r.rdb.Get(ctx, key).Result()

	if err != nil && err != redis.Nil {
		r.err = err
	}

	*value = val
	return r
}

func (r *Redis) Write(key string, value string) IRedis {
	if r.Err() != nil {
		return r
	}

	err := r.rdb.Set(ctx, key, value, 0).Err()

	if err != nil {
		r.err = err
	}

	return r
}

func (r *Redis) Delete(key string) IRedis {
	intCmd := r.rdb.Del(ctx, key)

	_, err := intCmd.Result()

	if err != nil {
		r.err = err
	}

	return r
}
