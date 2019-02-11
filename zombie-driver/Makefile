
TEST_HTTP_ADDR := 127.0.0.1:9001

.PHONY: all test

all:
	# Write build command(s) here

test:
	# Write test command(s) here
	@DRV_LOC_TEST_HTTP_ADDR=${TEST_HTTP_ADDR} \
		go test -race ${TARGS} ./...
