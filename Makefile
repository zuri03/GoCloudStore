resart-records: docker-clean
	docker rmi record_server -f \
	docker-compose up

restart-storage: docker-clean
	docker rmi gocloudstore_storage -f \
	docker-compose up

restart: docker-clean
	docker rmi gocloudstore_storage -f \
	docker rmi record_server -f \
	docker-compose up 

docker-clean:
	docker-compose down

