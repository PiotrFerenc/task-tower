DOCKER_COMPOSE_FILE := docker/docker-compose.yml

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
	docker build -t task-tower/controller -f docker/Dockerfile-controller .
	docker build -t task-tower/worker -f docker/Dockerfile-worker .
docker-rebuild:
	$(MAKE) docker-down
	$(MAKE) docker-build
	$(MAKE) docker-up
test:
	go test tests/main_test.go