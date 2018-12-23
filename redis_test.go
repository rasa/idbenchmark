package idbenchmark_test

import (
	"log"
	"testing"

	"github.com/go-redis/redis"
)

const (
	redisKey = idbenchmarkKey
)

func redisConnect() (db *redis.Client, err error) {
	db = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		PoolSize: 10,
	})

	_, err = db.Ping().Result()
	if err != nil {
		log.Printf("Sequence redis open error: %v", err)
		return nil, err
	}
	return db, nil
}

func runRedis(b *testing.B, db *redis.Client) {
	var id uint64

	for n := 0; n < b.N; n++ {
		err := db.Incr(redisKey).Err()
		if err != nil {
			log.Printf("Sequence redis incr error: %v", err)
			break
		}
		id, err = db.Get(redisKey).Uint64()
		if err != nil {
			log.Printf("Sequence redis get error: %v", err)
			break
		}
		if id == 0 {
			log.Printf("id=0")
			break
		}
	}
}

func BenchmarkRedis(b *testing.B) {
	db, err := redisConnect()
	if err != nil {
		return
	}
	defer db.Close()

	b.ResetTimer()
	runRedis(b, db)
	b.StopTimer()
}

func BenchmarkRedisParallel(b *testing.B) {
	db, err := redisConnect()
	if err != nil {
		return
	}
	defer db.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runRedis(b, db)
		}
	})
	b.StopTimer()
}
