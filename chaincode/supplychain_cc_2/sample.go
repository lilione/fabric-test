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

	if fn == "registerItem" {
		commitRegistrant := args[0]

		response := registerItem(stub, commitRegistrant)
		if response.Status == 200 {
			return shim.Success(response.Payload)
		} else {
			return shim.Error(response.Message)
		}

	} else if fn == "handOffItemToNextProvider" {
		commitInputProvider := args[0]
		commitOutputProvider := args[1]
		proof := args[2]
		itemID := args[3]
		prevSeq := args[4]
		fmt.Println("cc", args)

		response := handOffItemToNextProvider(stub, commitInputProvider, commitOutputProvider, proof, itemID, prevSeq)
		if response.Status == 200 {
			return shim.Success(response.Payload)
		} else {
			return shim.Error(response.Message)
		}

	}

	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
