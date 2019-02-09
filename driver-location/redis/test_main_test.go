package redis_test

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

var testRedisAddr string

func TestMain(m *testing.M) {
	testRedisAddr = os.Getenv("DRV_LOC_TEST_REDIS_ADDR")
	if testRedisAddr == "" {
		log.Fatal("DRV_LOC_TEST_REDIS_ADDR env var is required to run the test")
	}

	rand.Seed(time.Now().UnixNano())
	cleanRedis(testRedisAddr)
	os.Exit(m.Run())
}

func cleanRedis(addr string) {
	var c = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	var sc = c.FlushAll()
	if err := sc.Err(); err != nil {
		log.Fatalf("Connection to Redis failed or FlushAll command failed: %+v", err)
	}
}
