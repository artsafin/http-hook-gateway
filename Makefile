DEV_COMPOSE_FILE=deployments/dev/docker-compose.yml

test:
	docker build .

clean-dev:
	docker-compose -f ${DEV_COMPOSE_FILE} down --remove-orphans -v

run-dev:
	docker-compose -f ${DEV_COMPOSE_FILE} up -d --build --force-recreate
	docker-compose -f ${DEV_COMPOSE_FILE} logs -f
