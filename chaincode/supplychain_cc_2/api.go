package main

import "github.com/hyperledger/fabric-chaincode-go/shim"

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

func verify(stub shim.ChaincodeStubInterface, prevProvider string, sucProvider string, proofProvider string, prevAmt string, sucAmt string, proofAmt string) bool {
	chainCodeArgs := toChaincodeArgs("verify", prevProvider, sucProvider, proofProvider, prevAmt, sucAmt, proofAmt)
	response := stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
	if response.Status == 200 {
		return true
	}
	return false
}