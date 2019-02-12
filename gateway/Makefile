
BIN_NAME := gateway

.PHONY: all test clean

all:
	# Write build command(s) here
	@go build -o ${BIN_NAME} cmd/server/main.go

test:
	# Write test command(s) here
	@go test -race ./...

clean:
	@rm -f ${BIN_NAME}
