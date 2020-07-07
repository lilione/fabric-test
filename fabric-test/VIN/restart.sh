#!/usr/bin/env bash
#set -e

docker stop $(docker ps | grep 'client' | awk '{print $1}')
docker rm $(docker ps -a | grep "client" | awk '{print $1}')

# build system chaincode
echo "building system chaincode"
cd $GOPATH/src/github.com/lilione/fabric-test/chaincode/myscc
go build -buildmode=plugin

# install and instantiate application chaincode
cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
./byfn.sh restart -f docker-compose-from-docker.yaml -o kafka

# start mpc servers
echo "starting mpc servers"
cd $GOPATH/src/github.com/lilione/HoneyBadgerMPC
bash apps/fabric/scripts/start_server.sh

# start client
echo "staring client"
docker run -d \
 -v /Users/lilione/go/src/github.com/lilione/HoneyBadgerMPC:/usr/src/HoneyBadgerMPC \
 -v /Users/lilione/go/src/github.com/lilione/fabric-test/fabric-test/log/chaincode:/usr/src/HoneyBadgerMPC/apps/fabric/log/chaincode \
 -v /var/run/docker.sock:/var/run/docker.sock \
 --name client -it hyperledger/fabric-peer:latest
docker network connect net_byfn client
docker exec -it client bash

#docker logs peer0.org1.example.com
#docker logs peer0.org2.example.com
#docker logs peer1.org1.example.com
#docker logs peer1.org2.example.com