package http_test

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

var testHttpAddr string

func TestMain(m *testing.M) {
	testHttpAddr = os.Getenv("DRV_LOC_TEST_HTTP_ADDR")
	if testHttpAddr == "" {
		log.Fatal("DRV_LOC_TEST_HTTP_ADDR env var is required to run the test")
	}

	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}
