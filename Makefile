DOCKER_COMPOSE_FILE := config/docker-compose.yml

docker-reset:
	$(MAKE) docker-down
	$(MAKE) docker-up
docker-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
docker-down:config/docker-compose.yml
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
docker-clean:
	docker stop $(docker ps -aq)
	docker rm $(docker ps -aq)

