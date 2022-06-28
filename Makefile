build-recordKeeper:
	docker build . -f ./docker/dev/Dockerfile.record -t recordKeeper

build-storage:
	docker build . -f ./docker/dev/Dockerfile.storage -t storage

build-all: build-recordKeeper build-storage 

run: build-all
	docker run -d -p8000:8000 storage
	docker run -d -p8080:8080 recordKeeper

#broken target
clean:
	rm build/* build/.*; rmdir build
