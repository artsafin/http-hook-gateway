test:
	docker build .

run-dev:
	docker-compose down -f docker-compose.dev.yml --remove-orphans
	docker-compose up -f docker-compose.dev.yml --build
