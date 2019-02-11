package nsqhttp_test

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	testNSQdAddr       string
	testNSQLookupdAddr string
	testHttpAddr       string
)

func TestMain(m *testing.M) {
	testNSQdAddr = os.Getenv("DRV_LOC_TEST_NSQD_ADDR")
	if testNSQdAddr == "" {
		log.Fatal("DRV_LOC_TEST_NSQD_ADDR env var is required to run the test")
	}

	testNSQLookupdAddr = os.Getenv("DRV_LOC_TEST_NSQLOOKUPD_ADDR")
	if testNSQLookupdAddr == "" {
		log.Fatal("DRV_LOC_TEST_NSQLOOKUPD_ADDR env var is required to run the test")
	}

	testHttpAddr = os.Getenv("DRV_LOC_TEST_HTTP_ADDR")
	if testHttpAddr == "" {
		log.Fatal("DRV_LOC_TEST_HTTP_ADDR env var is required to run the test")
	}

	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}
