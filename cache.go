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
	GetCollectionProducts()
	// setCollectionProducts()
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

func setCollectionProducts(products []Products) {
	conn := Redis.Conn.Get()
	defer conn.Close()
	for _, v := range products {

		c, _ := json.Marshal(v.Products)

		_, err := conn.Do("SET", v.CategoryID, c)
		if err != nil {
			log.Println(err)

		}
	}
	// for _, p := range cp.Product {
	// 	setProducts(p.SKU)
	// }

}

func setProducts(sku string) {
	conn := Redis.Conn.Get()
	defer conn.Close()

	products := getProduct("d4a559dd-ea1e-40ea-951a-2c6effabf00a", sku)
	fmt.Println(products)
}

// Collections handles requests to get collection names and GUIDs
func Collections(w http.ResponseWriter, r *http.Request) {
	conn := Redis.Conn.Get()
	defer conn.Close()

	s, err := redis.Bytes(conn.Do("GET", "meta"))
	if err == redis.ErrNil {
		fmt.Fprintln(w, "Collection does not exist!")
	}

	col := []Meta{}

	err = json.Unmarshal([]byte(s), &col)
	collections, _ := json.Marshal(col)
	w.Header().Set("Content-Type", "application/json")
	w.Write(collections)
}

// CollectionProducts returns all product for a given collection
func CollectionProducts(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	// inputParams := strings.Split(r.URL.Path, "/")
	// collectionID := inputParams[2:]
	// fmt.Println(collectionID[0])
	id := r.FormValue("id")

	conn := Redis.Conn.Get()
	defer conn.Close()

	s, err := redis.Bytes(conn.Do("GET", id))
	if err == redis.ErrNil {
		fmt.Fprintln(w, "The Product for that collection does not exist!")
		return
	}

	var p []Product

	err = json.Unmarshal([]byte(s), &p)
	products, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(products)

}
