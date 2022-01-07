package message

import (
	"github.com/wearedevx/keystone/api/internal/redis"
)

// var Redis *redis.Redis

type messageService struct {
	redis *redis.Redis
}

type MessageService interface {
	GetMessageByUuid(uuid string) ([]byte, error)
	WriteMessageWithUuid(uuid string, value []byte) error
	DeleteMessageWithUuid(uuid string) error
}

func NewMessageService() MessageService {
	return &messageService{
		redis: redis.NewRedis(),
	}
}

func (m *messageService) GetMessageByUuid(uuid string) ([]byte, error) {
	value := ""
	m.redis.Read(uuid, &value)

	if m.redis.Err() != nil {
		return nil, m.redis.Err()
	}

	return []byte(value), nil
}

func (m *messageService) WriteMessageWithUuid(uuid string, value []byte) error {
	m.redis.Write(uuid, string(value))
	return m.redis.Err()
}

func (m *messageService) DeleteMessageWithUuid(uuid string) error {
	m.redis.Delete(uuid)

	return m.redis.Err()
}
