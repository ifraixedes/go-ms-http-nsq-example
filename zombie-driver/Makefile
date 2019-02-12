
BIN_NAME := zombiedrv
TEST_HTTP_ADDR := 127.0.0.1:9001

.PHONY: all test clean

all:
	# Write build command(s) here
	@go build -o ${BIN_NAME} ./http/cmd/server/main.go

test:
	# Write test command(s) here
	@DRV_LOC_TEST_HTTP_ADDR=${TEST_HTTP_ADDR} \
		go test -race ${TARGS} ./...

clean:
	@rm -f ${BIN_NAME}
