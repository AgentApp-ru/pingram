package redis

import (
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

func GetInt64(key string) (int64, error) {
	conn := Redis.Get()
	defer conn.Close()

	var (
		data int64
		err  error
	)
	data, err = redigo.Int64(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}

	return data, nil
}

func GetBytes(key string) ([]byte, error) {
	conn := Redis.Get()
	defer conn.Close()

	var (
		data []byte
		err  error
	)
	data, err = redigo.Bytes(conn.Do("GET", key))
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func SetBytes(key string, value []byte) error {
	conn := Redis.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func CreateOrUpdateFull(key string, value int64) error {
	conn := Redis.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func CreateOrUpdateShort(key string) error {
	return CreateOrUpdateFull(key, time.Now().Unix())
}

func Del(key string) error {
	conn := Redis.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func GetKeys(keysPattern string) ([]string, error) {
	conn := Redis.Get()
	defer conn.Close()

	keys, err := redigo.Strings(conn.Do("KEYS", keysPattern))
	if err != nil {
		return []string{}, err
	}

	return keys, nil
}
