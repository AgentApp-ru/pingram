package redis

import (
	"pingram/pkg/config"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var (
	Redis *redigo.Pool
)

func init() {
	Redis = GetRedisPool(config.Settings.RedisUrl)
}

func GetRedisPool(urls string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     10,
		IdleTimeout: 180 * time.Second, // Default is 300 seconds for redis server
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", urls)
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}
