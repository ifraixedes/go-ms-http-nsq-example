.PHONY: all test dev-env clean e2e

all:
	make -C ./driver-location
	make -C ./gateway
	make -C ./zombie-driver

test:
	make -C ./driver-location test
	make -C ./gateway test
	make -C ./zombie-driver test

dev-env:
	@docker-compose up -d

clean:
	@docker-compose down
	make -C ./driver-location clean
	make -C ./gateway clean
	make -C ./zombie-driver clean

## This test works with the current docker-compose, and running the 3 binaries
## (gateway, zombiedrv and driverloc) with the default configuration and the
## current gateway/config.yaml file is used.
e2e:
	@curl -v -X PATCH -H "Content-Type: application/json" \
		-d '{"latitude": 41.388038, "longitude": 2.18249}' 127.0.0.1:9002/drivers/1/locations
	@curl -v -X PATCH -H "Content-Type: application/json" \
		-d '{"latitude": 41.402469, "longitude": 2.187417}' 127.0.0.1:9002/drivers/1/locations
	@curl -v -X GET 127.0.0.1:9002/drivers/1
	@curl -v -X PATCH -H "Content-Type: application/json" \
		-d '{"latitude": 41.388038, "longitude": 2.18249}' 127.0.0.1:9002/drivers/2/locations
	@curl -v -X PATCH -H "Content-Type: application/json" \
		-d '{"latitude": 41.388038, "longitude": 2.18249}' 127.0.0.1:9002/drivers/2/locations
	@curl -v -X GET 127.0.0.1:9002/drivers/2
