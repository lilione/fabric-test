package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
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

func verifyEq(stub shim.ChaincodeStubInterface, commitPrev string, commitSuc string, proof string) bool {
	chainCodeArgs := toChaincodeArgs("verifyEq", commitPrev, commitSuc, proof)
	response := stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
	if response.Status == 200 {
		return true
	}
	return false
}