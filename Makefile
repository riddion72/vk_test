all: build run

build:
	docker build -t container2 ./some_container2
	docker-compose build --no-cache

run:
	docker run --name container2 -d -t -i container2
	docker-compose up -d

restart:
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up -d

stop:
	docker-compose down -v
	docker stop container2