#!/usr/bin/env bash

#docker rm all containers
docker rm -f $(docker ps -aq)

docker ps -a

#clearup docker volumes
docker system prune --volumes

docker volume ls

docker rmi setup-honeybadger