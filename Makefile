.PHONY: test build

test:
	go test ./...

clean:
	go clean
	-rm -r ./build/*

build:
	make clean
	go build -o build ./cmd/search
	go build -o build ./cmd/settings
	cp ./configs/* build/
	cp -r ./assets/* build/
