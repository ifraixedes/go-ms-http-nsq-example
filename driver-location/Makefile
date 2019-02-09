
TEST_REDIS_ADDR := 127.0.0.1:6379

.PHONY: all test clean

all:
	# Write your build command(s) here

test:
# TARGS are the CLI go test flags which allows to pass any additional falg when
# executing this target.
	@DRV_LOC_TEST_REDIS_ADDR=${TEST_REDIS_ADDR} go test -race ${TARGS} ./...
