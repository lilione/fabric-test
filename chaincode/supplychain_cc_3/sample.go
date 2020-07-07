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

	} else if fn == "createTruck" {
		response := createTruck(stub)
		return shim.Success(response.Payload)
	} else if fn == "recordShipment" {
		truckID := args[0]
		idxLoadTime := args[1]
		maskedLoadTime := args[2]
		idxUnloadTime := args[3]
		maskedUnloadTime := args[4]

		response := recordShipment(stub, truckID, idxLoadTime, maskedLoadTime, idxUnloadTime, maskedUnloadTime)
		if response.Status == 200 {
			return shim.Success(response.Payload)
		} else {
			return shim.Error(response.Message)
		}
	} else if fn == "queryPositions" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]

		response := queryPositions(stub, truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime)
		return shim.Success(response.Payload)
	} else if fn == "queryNumber" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]

		response := queryNumber(stub, truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime)
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
