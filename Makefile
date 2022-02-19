.PHONY: test build

test:
	go test ./...

clean:
	go clean
	-rm ./build/list
	-rm ./build/settings

build:
	make clean
	go build -o build ./cmd/list
	go build -o build ./cmd/settings
