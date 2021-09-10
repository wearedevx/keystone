package message

import (
	. "github.com/wearedevx/keystone/api/internal/redis"
)

// var Redis *redis.Redis

type MessageService struct {
	redis *Redis
}

func NewMessageService() *MessageService {
	return &MessageService{
		redis: NewRedis(),
	}
}

func (m *MessageService) GetMessageByUuid(uuid string) ([]byte, error) {
	value := ""
	m.redis.Read(uuid, &value)

	if m.redis.Err() != nil {
		return nil, m.redis.Err()
	}

	m.redis.Delete(uuid)

	return []byte(value), nil
}

func (m *MessageService) WriteMessageWithUuid(uuid string, value []byte) error {
	m.redis.Write(uuid, string(value))
	return m.redis.Err()
}
