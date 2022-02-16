.PHONY: test build

test:
	go test ./...

clean:
	go clean
	-rm ./build/list
	-rm ./build/config

build:
	make clean
	go build -o build ./cmd/list
	go build -o build ./cmd/config
