.PHONY: test build

test:
	go test ./...

clean:
	go clean
	-rm ./build/search
	-rm ./build/settings

build:
	make clean
	go build -o build ./cmd/search
	go build -o build ./cmd/settings
