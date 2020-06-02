#!/usr/bin/env bash
set -e

cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test/VIN

docker volume create orderer.example.com-crypto
docker volume create peer0.org1.example.com-msp
docker volume create peer0.org2.example.com-msp
docker volume create peer1.org1.example.com-msp
docker volume create peer1.org2.example.com-msp
docker volume create peer0.org1.example.com-tls
docker volume create peer0.org2.example.com-tls
docker volume create peer1.org1.example.com-tls
docker volume create peer1.org2.example.com-tls

# build docker image for setup
#docker build -t setup-honeybadger .

# run docker container
docker run -it \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ${PWD}/setup.sh:/opt/setup.sh \
	-v ${PWD}/restart.sh:/opt/restart.sh \
	-v $GOPATH/src/github.com/lilione/HoneyBadgerMPC:/opt/gopath/src/github.com/lilione/HoneyBadgerMPC \
  -v $GOPATH/src/github.com/lilione/fabric-test:/opt/gopath/src/github.com/lilione/fabric-test \
  -v $GOPATH/src/github.com/lilione/fabric:/opt/gopath/src/github.com/hyperledger/fabric \
	-v orderer.example.com-crypto:/opt/crypto/orderer.example.com \
	-v peer0.org1.example.com-tls:/opt/crypto/peer0.org1.example.com/tls \
	-v peer0.org1.example.com-msp:/opt/crypto/peer0.org1.example.com/msp \
	-v peer1.org1.example.com-tls:/opt/crypto/peer1.org1.example.com/tls \
	-v peer1.org1.example.com-msp:/opt/crypto/peer1.org1.example.com/msp \
	-v peer0.org2.example.com-tls:/opt/crypto/peer0.org2.example.com/tls \
	-v peer0.org2.example.com-msp:/opt/crypto/peer0.org2.example.com/msp \
	-v peer1.org2.example.com-tls:/opt/crypto/peer1.org2.example.com/tls \
	-v peer1.org2.example.com-msp:/opt/crypto/peer1.org2.example.com/msp \
	--name setup-honeybadger \
	-w /opt setup-honeybadger \
	bash

# enter stopped container
#docker restart setup-honeybadger

# enter running container
#docker exec -it setup-honeybadger bash