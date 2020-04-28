#!/usr/bin/env bash

cd ..
./byfn.sh down -f docker-compose-from-docker.yaml -o kafka

#docker rm all containers
docker stop $(docker ps -aq)
docker rm -f $(docker ps -aq)
docker ps -a

#clearup docker volumes
docker system prune --volumes
docker volume ls

#docker rmi setup-honeybadger