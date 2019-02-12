
BIN_NAME := driverloc
TEST_REDIS_ADDR := 127.0.0.1:6379
TEST_NSQD_ADDR := 127.0.0.1:4150
TEST_NSQLOOKUPD_ADDR := 127.0.0.1:4161
TEST_HTTP_ADDR := 127.0.0.1:9000

.PHONY: all test clean

all:
	@go build -o ${BIN_NAME} ./nsqhttp/cmd/server/main.go

test:
# TARGS are the CLI go test flags which allows to pass any additional falg when
# executing this target.
	@DRV_LOC_TEST_REDIS_ADDR=${TEST_REDIS_ADDR} \
		DRV_LOC_TEST_NSQD_ADDR=${TEST_NSQD_ADDR} \
		DRV_LOC_TEST_NSQLOOKUPD_ADDR=${TEST_NSQLOOKUPD_ADDR} \
		DRV_LOC_TEST_HTTP_ADDR=${TEST_HTTP_ADDR} \
		go test -race ${TARGS} ./...

clean:
	@rm -f ${BIN_NAME}
