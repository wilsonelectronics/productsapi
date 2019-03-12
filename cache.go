package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisHelper struct {
	Conn *redis.Pool
}

// Redis . . .
var Redis *redisHelper

func init() {
	fmt.Println("Cache.go hit!!")
	Redis = &redisHelper{Conn: newPool()}
	setCollections()
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			return redis.DialURL(os.Getenv("REDIS_URL"))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := redis.String(c.Do("PING"))
			return err
		},
	}
}

func setCollections() {
	conn := Redis.Conn.Get()
	defer conn.Close()

	collection := GetCollections()

	json, err := json.Marshal(collection)

	_, err = conn.Do("SET", "meta", json)
	if err != nil {
		log.Println(err)
	}
}

// GetCollection handles requests to get collection names and GUIDs
func GetCollection(w http.ResponseWriter, r *http.Request) {
	conn := Redis.Conn.Get()
	defer conn.Close()

	s, err := redis.Bytes(conn.Do("GET", "meta"))
	if err == redis.ErrNil {
		log.Println("Collection does not exist!")
	}

	col := []Meta{}

	err = json.Unmarshal([]byte(s), &col)
	collections, _ := json.Marshal(col)
	w.Header().Set("Content-Type", "application/json")
	w.Write(collections)
}
