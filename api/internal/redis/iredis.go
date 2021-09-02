package redis

type IRedis interface {
	Read(key string, value *string) IRedis
	Write(key string, value string) IRedis
	Delete(key string) IRedis
}
