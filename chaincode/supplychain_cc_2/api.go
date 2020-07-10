package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	channelName = "mychannel"
	sccName = "supplychain_scc_2"
)

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func registerItem(stub shim.ChaincodeStubInterface, commitRegistrant string) peer.Response {
	chainCodeArgs := toChaincodeArgs("registerItem", commitRegistrant)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func handOffItemToNextProvider(stub shim.ChaincodeStubInterface, commitInputProvider string, commitOutputProvider string, proof string, itemID string, prevSeq string) peer.Response {
	chainCodeArgs := toChaincodeArgs("handOffItemToNextProvider", commitInputProvider, commitOutputProvider, proof, itemID, prevSeq)
	fmt.Println("chainCodeArgs", chainCodeArgs)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}