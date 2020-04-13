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
	function, args := stub.GetFunctionAndParameters()

	//if function == "query" {
	//	response := queryScc(stub, args[0])
	//	return shim.Success(response.Payload)
	//} else if function == "update" {
	//	response := updateScc(stub, args[0], args[1])
	//	return shim.Success(response.Payload)
	//}
	if function == "query" {
		key := args[0]
		value, err := stub.GetState(key)
		if err != nil {
			return shim.Error(fmt.Sprint("failed to get value for key %s", key))
		}
		return shim.Success([]byte(value))
	} else if function == "update" {
		key, value := args[0], []byte(args[1])
		err := stub.PutState(key, value)
		if err != nil {
			return shim.Error(fmt.Sprint("failed to update for key %s value %s", key, value))
		}
		return shim.Success(nil)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
