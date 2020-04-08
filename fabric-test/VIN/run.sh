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
docker volume create honeybadgerscc
docker volume create chaincode
docker volume create cli-config
docker volume create peer-bin

#build docker image for setup
docker build -t setup-honeybadger .

#run docker container
docker run -it \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ${PWD}/setup.sh:/opt/setup.sh \
	-v peer-bin:/opt/peer-bin \
	-v orderer.example.com-crypto:/opt/crypto/orderer.example.com \
	-v peer0.org1.example.com-tls:/opt/crypto/peer0.org1.example.com/tls \
	-v peer0.org1.example.com-msp:/opt/crypto/peer0.org1.example.com/msp \
	-v peer1.org1.example.com-tls:/opt/crypto/peer1.org1.example.com/tls \
	-v peer1.org1.example.com-msp:/opt/crypto/peer1.org1.example.com/msp \
	-v peer0.org2.example.com-tls:/opt/crypto/peer0.org2.example.com/tls \
	-v peer0.org2.example.com-msp:/opt/crypto/peer0.org2.example.com/msp \
	-v peer1.org2.example.com-tls:/opt/crypto/peer1.org2.example.com/tls \
	-v peer1.org2.example.com-msp:/opt/crypto/peer1.org2.example.com/msp \
	-v honeybadgerscc:/opt/gopath/src/github.com/lilione/honeybadgerscc \
	-v chaincode:/opt/chaincode \
	-v cli-config:/opt/crypto/cli \
	--name setup-honeybadger \
	-w /opt setup-honeybadger \
	bash

#enter stopped container
#docker restart setup-honeybadger

#enter running container
#docker exec -it setup-honeybadger bash


