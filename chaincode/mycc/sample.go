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
		num := args[0]
		response := getInputmaskIdx(stub, num)
		return shim.Success(response.Payload)

	} else if fn == "registerItem" {
		idxRegistrant := args[0]
		maskedRegistrant := args[1]
		idxAmt := args[2]
		maskedAmt := args[3]
		response := registerItem(stub, idxRegistrant, maskedRegistrant, idxAmt, maskedAmt)
		return shim.Success(response.Payload)

	} else if fn == "handOffItemToNextProvider" {
		idxInputProvider := args[0]
		maskedInputProvider := args[1]
		idxOutputProvider := args[2]
		maskedOutputProvider := args[3]
		idxAmt := args[4]
		maskedAmt := args[5]
		itemID := args[6]
		prevSeq := args[7]
		response := handOffItemToNextProvider(stub, idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq)
		return shim.Success(response.Payload)

	} else if fn == "sourceItem" {
		itemID := args[0]
		seq := args[1]
		response := sourceItem(stub, itemID, seq)
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
