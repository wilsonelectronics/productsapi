package cache

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool = &redis.Pool{
	MaxIdle:     10,
	IdleTimeout: 240 * time.Second,

	Dial: func() (redis.Conn, error) {
		return redis.DialURL(os.Getenv("REDIS_URL"))
	},
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := redis.String(c.Do("PING"))
		return err
	}}

// Retrieve . . .
func Retrieve(key string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return redis.Bytes(conn.Do("GET", key))
}

// Store . . .
func Store(key string, bytes []byte) error {
	conn := pool.Get()
	defer conn.Close()

	secondsInDay := 60
	_, err := conn.Do("SETEX", key, secondsInDay, bytes)
	return err
}
