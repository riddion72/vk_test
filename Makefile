all:
	docker-compose build --no-cache
	docker compose up -d

restart:
	docker compose down -v
	docker-compose build --no-cache
	docker compose up -d

stop:
	docker compose down -v