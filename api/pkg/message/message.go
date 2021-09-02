package message

import (
	"github.com/wearedevx/keystone/api/internal/redis"
)

var Redis *redis.Redis

func GetMessageByUuid(uuid string) ([]byte, error) {
	value := ""
	Redis.Read(uuid, &value)

	if Redis.Err() != nil {
		return nil, Redis.Err()
	}

	Redis.Delete(uuid)

	return []byte(value), nil
}

func WriteMessageWithUuid(uuid string, value []byte) error {
	Redis.Write(uuid, string(value))
	return Redis.Err()
}

func init() {
	Redis = new(redis.Redis)
}
