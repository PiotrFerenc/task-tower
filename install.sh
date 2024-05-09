docker ps -q --filter "ancestor=task-tower/controller" | grep -q . && docker stop $(docker ps -q --filter "ancestor=task-tower/controller") || echo "No container to stop"
docker ps -q --filter "ancestor=task-tower/worker" | grep -q . && docker stop $(docker ps -q --filter "ancestor=task-tower/worker") || echo "No container to stop"
docker ps -aq --filter "ancestor=task-tower/controller" | grep -q . && docker rm -f $(docker ps -aq --filter "ancestor=task-tower/controller") || echo "No container to remove"
docker ps -aq --filter "ancestor=task-tower/worker" | grep -q . && docker rm -f $(docker ps -aq --filter "ancestor=task-tower/worker") || echo "No container to remove"
