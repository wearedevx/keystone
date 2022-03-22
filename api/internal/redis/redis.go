// +build !test

package redis

import (
	"context"
	"math"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	err error
	rdb *redis.Client
}

var rdb *redis.Client

var (
	redisHost  string
	redisPort  string
	redisIndex string
)

var ctx = context.Background()

func getOrDefault(value string, envkey string, defaultValue string) string {
	if value != "" {
		return value
	}

	if value = os.Getenv(envkey); value != "" {
		return value
	}

	return defaultValue
}

func NewRedis() *Redis {
	var err error
	var r Redis

	redisHost = getOrDefault(redisHost, "REDIS_HOST", "localhost")
	redisPort = getOrDefault(redisPort, "REDIS_PORT", "6379")
	redisIndex = getOrDefault(redisIndex, "REDIS_INDEX", "0")

	redisIndexInt, err := strconv.ParseInt(redisIndex, 10, 64)
	if err != nil {
		redisIndexInt = 0
	}

	// FIXME: shouldn’t we log this ?
	if redisIndexInt > math.MaxInt {
		redisIndexInt = 0
	}

	r.rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",                 // no password set
		DB:       int(redisIndexInt), // use default DB
	})

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

// For tests only, this is a no-op
func (r *Redis) SetupFixtures(_ map[string]string) {
}
