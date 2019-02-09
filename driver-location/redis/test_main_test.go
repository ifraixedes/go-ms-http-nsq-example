package redis_test

import (
	"log"
	"os"
	"testing"
)

var testRedisAddr string

func TestMain(m *testing.M) {
	testRedisAddr = os.Getenv("DRV_LOC_TEST_REDIS_ADDR")
	if testRedisAddr == "" {
		log.Fatal("DRV_LOC_TEST_REDIS_ADDR env var is required to run the test")
	}
	os.Exit(m.Run())
}
