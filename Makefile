docker-build-recordKeeper:
	docker build . -f ./deployments/dev/Dockerfile.records -t records

docker-build-storage:
	docker build . -f ./deployments/dev/Dockerfile.storage -t storage

docker-build-all: build-recordKeeper build-storage 

docker-run: build-all
	docker run -d -p8000:8000 storage
	docker run -d -p8080:8080 records
	cd cmd/cli && go run main.go

