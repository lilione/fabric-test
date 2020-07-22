#!/usr/bin/env bash
set -e

docker stop $(docker ps | grep 'client' | awk '{print $1}') || true
docker rm $(docker ps -a | grep "client" | awk '{print $1}') || true

docker rmi -f $(docker images | grep 'dev' | awk '{print $1}') || true

rm -rfv $GOPATH/src/github.com/lilione/fabric-test/fabric-test/log/chaincode

## build system chaincode
#echo "building system chaincode"
#cd $GOPATH/src/github.com/lilione/fabric-test/chaincode/supplychain_scc_1
#go build -buildmode=plugin
#cd $GOPATH/src/github.com/lilione/fabric-test/chaincode/supplychain_scc_2
#go build -buildmode=plugin
#cd $GOPATH/src/github.com/lilione/fabric-test/chaincode/supplychain_scc_3
#go build -buildmode=plugin

## build fabric peer binary
#echo "building fabric-peer binary"
#cd $GOPATH/src/github.com/hyperledger/fabric
#GO_TAGS+=" pluginsenabled" make peer

# build fabric peer docker image
#echo "building fabric-peer docker image"
#DOCKER_DYNAMIC_LINK=true GO_TAGS+=" pluginsenabled" make peer-docker IN_DOCKER=true

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

# docker logs peer0.org1.example.com
# docker logs peer0.org2.example.com
# docker logs peer1.org1.example.com
# docker logs peer1.org2.example.com
# docker logs cli