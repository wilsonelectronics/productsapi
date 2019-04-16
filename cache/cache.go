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

	_, err := conn.Do("SET", key, bytes)
	return err
}

// func setCollections() {
// 	collection := GetCollections()
// 	json, err := json.Marshal(collection)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	_, err = conn.Do("SET", "meta", json)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func setCollectionProducts() {
// 	products := GetCollectionProducts()
// 	conn := Redis.Conn.Get()
// 	defer conn.Close()
// 	for _, v := range products {

// 		c, _ := json.Marshal(v.Products)

// 		_, err := conn.Do("SET", v.CategoryID, c)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

// // SetSingleProduct attempts to retrieve a single product from the cache
// func SetSingleProduct(product []byte) error {
// 	productJSON, err := json.Marshal(product)
// 	if err != nil {
// 		return err
// 	}

// 	conn := Redis.Conn.Get()
// 	defer conn.Close()
// 	_, err = conn.Do("SET", product.ProductSiteGUID, productJSON)
// 	return err
// }

// // GetSingleProduct attempts to retrieve a single product from the cache, returns a nil product with no error if not found
// func GetSingleProduct(productSiteGUID string) ([]byte, error) {
// 	conn := Redis.Conn.Get()
// 	defer conn.Close()

// 	s, err := redis.Bytes(conn.Do("GET", productSiteGUID))
// 	if err == redis.ErrNil {
// 		return nil, nil
// 	}
// 	return s, err
// }

// // GetCachedCollections handles requests to get collection names and GUIDs
// func GetCachedCollections(w http.ResponseWriter, r *http.Request) {
// 	conn := Redis.Conn.Get()
// 	defer conn.Close()

// 	s, err := redis.Bytes(conn.Do("GET", "meta"))
// 	if err == redis.ErrNil {
// 		fmt.Fprintln(w, "Collection does not exist!")
// 	}

// 	col := []Meta{}

// 	err = json.Unmarshal([]byte(s), &col)
// 	collections, _ := json.Marshal(col)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(collections)
// }

// // GetCachedCollectionProducts returns all product for a given collection
// func GetCachedCollectionProducts(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL.Path)
// 	id := r.FormValue("id")

// 	conn := Redis.Conn.Get()
// 	defer conn.Close()

// 	s, err := redis.Bytes(conn.Do("GET", id))
// 	if err == redis.ErrNil {
// 		fmt.Fprintln(w, "The Product for that collection does not exist!")
// 		return
// 	}

// 	var p []Product

// 	err = json.Unmarshal([]byte(s), &p)
// 	products, _ := json.Marshal(p)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(products)
// }
