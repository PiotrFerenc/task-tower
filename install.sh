#!/usr/bin/zsh
controller_container="task-tower/controller"
worker_container="task-tower/worker"

docker ps -q --filter "ancestor=$controller_container" | grep -q . && docker stop $(docker ps -q --filter "ancestor=$controller_container") || echo "No container to stop"
docker ps -q --filter "ancestor=$worker_container" | grep -q . && docker stop $(docker ps -q --filter "ancestor=$controller_container") || echo "No container to stop"

docker ps -aq --filter "ancestor=$worker_container" | grep -q . && docker rm -f $(docker ps -aq --filter "ancestor=$worker_container") || echo "No container to remove"
docker ps -aq --filter "ancestor=$controller_container" | grep -q . && docker rm -f $(docker ps -aq --filter "ancestor=$worker_container") || echo "No container to remove"
