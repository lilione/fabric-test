package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type SmartContract struct {
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "getInputmaskIdx" {
		response := getInputmaskIdx(stub)
		return shim.Success(response.Payload)
	} else if fn == "sendMaskedInput" {
		idx := args[0]
		maskedInput := args[1]

		response := sendMaskedInput(stub, idx, maskedInput)
		return shim.Success(response.Payload)
	} else if fn == "reconstruct" {
		idx := args[0]
		response := reconstruct(stub, idx)
		return shim.Success(response.Payload)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
