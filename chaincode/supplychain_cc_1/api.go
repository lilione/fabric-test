package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

const (
	channelName = "mychannel"
	sccName = "supplychain_scc_1"
)

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func registerItem(stub shim.ChaincodeStubInterface, args string) {
	chainCodeArgs := toChaincodeArgs("registerItem", args)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func handOffItem(stub shim.ChaincodeStubInterface, args string) {
	chainCodeArgs := toChaincodeArgs("handOffItem", args)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func sourceItem(stub shim.ChaincodeStubInterface, args string) {
	chainCodeArgs := toChaincodeArgs("sourceItem", args)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}