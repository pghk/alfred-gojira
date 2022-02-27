.PHONY: test build

test:
	go test ./...

clean:
	go clean
	-rm -r ./build/*

build:
	make clean
	go build -o build ./cmd/list
	go build -o build ./cmd/configure
	cp ./configs/* build/
	cp -r ./assets/* build/
