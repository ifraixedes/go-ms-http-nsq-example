.PHONY: all test dev-env clean

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
