all: build

build:
	go build ./cmd/endpoint_server/

clean:
	rm -f endpoint_server
