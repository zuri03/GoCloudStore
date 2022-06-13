build-cli:
	go build -o build/cli cmd/client/main.go

build-recordKeeper:
	go build -o build/recordKeeper cmd/recordKeeper/main.go

build-storage:
	go build -o build/storage cmd/storage/main.go

build-all: build-cli build-recordKeeper build-storage

run: build-all
	go run build/recordKeeper
	go run build/storage
	go run build/cli cli

#broken target
clean:
	rm build/* build/.*; rmdir build
