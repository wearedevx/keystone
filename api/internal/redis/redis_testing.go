// +build test

package redis

type Redis struct {
	err error
}

var fakeRedis map[string]string

func NewRedis() *Redis {
	fakeRedis = make(map[string]string)

	return new(Redis)
}

func (r *Redis) Err() error {
	return r.err
}

func (r *Redis) Read(key string, value *string) IRedis {
	*value = fakeRedis[key]
	return r
}

func (r *Redis) Write(key string, value string) IRedis {
	fakeRedis[key] = value
	return r
}

func (r *Redis) Delete(key string) IRedis {
	delete(fakeRedis, key)
	return r
}
