$controller_container = "task-tower/controller"
$worker_container = "task-tower/worker"

# Stop the controller container if it is running
$controllerIDs = docker ps -q --filter "ancestor=$controller_container"
if ($controllerIDs) {
    docker stop $controllerIDs
} else {
    Write-Host "No controller container to stop"
}

# Stop the worker container if it is running
$workerIDs = docker ps -q --filter "ancestor=$worker_container"
if ($workerIDs) {
    docker stop $workerIDs
} else {
    Write-Host "No worker container to stop"
}

# Remove the worker container if it exists
$workerAllIDs = docker ps -aq --filter "ancestor=$worker_container"
if ($workerAllIDs) {
    docker rm -f $workerAllIDs
} else {
    Write-Host "No worker container to remove"
}

# Remove the controller container if it exists
$controllerAllIDs = docker ps -aq --filter "ancestor=$controller_container"
if ($controllerAllIDs) {
    docker rm -f $controllerAllIDs
} else {
    Write-Host "No controller container to remove"
}
