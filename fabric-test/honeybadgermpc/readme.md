Build peer and peer docker images and run install and instantiate rockpaperscissors chaincode

```
docker volume create orderer.example.com-crypto 
docker volume create peer0.org1.example.com-crypto 
docker volume create peer0.org2.example.com-crypto 
docker volume create peer1.org1.example.com-crypto 
docker volume create peer1.org2.example.com-crypto 
docker volume create honeybadgerscc 
docker volume create chaincode 
docker volume create cli-config
docker volume create peer-bin
docker run -it \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ${PWD}/setup.sh:/opt/setup.sh \
	-v ${PWD}/docker-compose-from-docker.yaml:/opt/docker-compose-from-docker.yaml \
	-v peer-bin:/opt/peer-bin \
	-v orderer.example.com-crypto:/opt/crypto/orderer.example.com \
	-v peer0.org1.example.com-crypto:/opt/crypto/peer0.org1.example.com \
	-v peer1.org1.example.com-crypto:/opt/crypto/peer1.org1.example.com \
	-v peer0.org2.example.com-crypto:/opt/crypto/peer0.org2.example.com \
	-v peer1.org2.example.com-crypto:/opt/crypto/peer1.org2.example.com \
	-v honeybadgerscc:/opt/gopath/src/github.com/njeans/honeybadgerscc \
	-v chaincode:/opt/chaincode \
	-v cli-config:/opt/crypto/cli \
	--name setup-honeybadger \
	-w /opt hyperledger/fabric-baseimage:amd64-0.4.14 \
	bash setup.sh
```

Cleanup

```
docker container rm setup-honeybader
docker volume rm orderer.example.com-crypto peer0.org1.example.com-crypto peer0.org2.example.com-crypto peer1.org1.example.com-crypto peer1.org2.example.com-crypto honeybadgerscc chaincode cli-config peer-bin
```
