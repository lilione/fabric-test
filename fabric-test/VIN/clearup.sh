#!/usr/bin/env bash

docker stop client
docker rm client

cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
rm mychannel.block
rm -rf log
rm log.txt

# docker rm all containers
docker stop $(docker ps -aq)
docker rm -f $(docker ps -aq)
docker ps -a

# clearup docker volumes
docker system prune --volumes
docker volume ls

# clearup docker images
docker rmi $(docker images | grep 'dev')
#docker rmi setup-honeybadger