package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const channelName = "mychannel"

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func queryMyscc(stub shim.ChaincodeStubInterface, key string) peer.Response {
	chainCodeArgs := toChaincodeArgs("query", key)
	return stub.InvokeChaincode("myscc", chainCodeArgs, channelName)
}

func updateMyscc(stub shim.ChaincodeStubInterface, key string, value string) peer.Response {
	chainCodeArgs := toChaincodeArgs("update", key, value)
	return stub.InvokeChaincode("myscc", chainCodeArgs, channelName)
}
