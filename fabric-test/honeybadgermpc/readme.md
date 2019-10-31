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
docker build -t setup-honeybadger .
docker run -it \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ${PWD}/setup.sh:/opt/setup.sh \
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
	-w /opt setup-honeybadger \
	bash setup.sh
```


Restart network from fabric-test/fabric-test directory
```
cd ..
mkdir crypto-config
./byfn.sh restart -f docker-compose-from-docker.yaml
```

Cleanup (from fabric-test/fabric-test directory)

```
cd ..
./byfn.sh down -f docker-compose-from-docker.yaml
docker container rm setup-honeybadger
docker volume rm orderer.example.com-crypto peer0.org1.example.com-crypto peer0.org2.example.com-crypto peer1.org1.example.com-crypto peer1.org2.example.com-crypto honeybadgerscc chaincode cli-config peer-bin
```
