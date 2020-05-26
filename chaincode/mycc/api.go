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

func getInputmaskIdx(stub shim.ChaincodeStubInterface) peer.Response {
	chainCodeArgs := toChaincodeArgs("getInputmaskIdx")
	return stub.InvokeChaincode("myscc", chainCodeArgs, channelName)
}

func sendMaskedInput(stub shim.ChaincodeStubInterface, idx string, maskedInput string) peer.Response {
	chainCodeArgs := toChaincodeArgs("sendMaskedInput", idx, maskedInput)
	return stub.InvokeChaincode("myscc", chainCodeArgs, channelName)
}

func reconstruct(stub shim.ChaincodeStubInterface, idx string) peer.Response {
	chainCodeArgs := toChaincodeArgs("reconstruct", idx)
	return stub.InvokeChaincode("myscc", chainCodeArgs, channelName)
}