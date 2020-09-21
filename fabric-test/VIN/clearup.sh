#!/usr/bin/env bash

docker stop client
docker rm client

# docker rm all containers
docker stop $(docker ps -aq)
docker rm -f $(docker ps -aq)
docker ps -a

# clearup docker volumes
y | docker system prune --volumes
docker volume ls

# clearup docker images
docker rmi $(docker images | grep 'dev' | awk '{print $1}')
#docker rmi setup-honeybadger

cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
rm mychannel.block
rm -rf log
rm log.txt

cd $GOPATH/src/github.com/lilione/HoneyBadgerMPC/apps/fabric/log/server
rm log_0.txt
rm log_1.txt
rm log_2.txt
rm log_3.txt