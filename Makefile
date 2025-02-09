all: build run

build:
	docker-compose build --no-cache
	docker build -t container2 ./some_container2

run:
	docker-compose up -d
	docker run -d --name container2 -p 9876:80 --rm container2

restart:
	docker-compose down -v
	docker rm -f container2 || true
	docker-compose build --no-cache
	docker-compose up -d

stop:
	docker-compose down -v
	docker stop container2 || true

