DOCKER_COMPOSE_FILE := config/docker-compose.yml

docker-reset:
	$(MAKE) docker-down
	$(MAKE) docker-up
docker-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --remove-orphans
docker-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
docker-clean:
	docker stop $(docker ps -aq)
	docker rm $(docker ps -aq)

docker-build:
	docker build -t dwas/controller -f config/controller/Dockerfile .
	docker build -t dwas/worker -f config/worker/Dockerfile .
docker-rebuild:
	$(MAKE) docker-down
	$(MAKE) docker-build
	$(MAKE) docker-up
