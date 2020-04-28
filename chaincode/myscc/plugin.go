/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

// New returns an implementation of the chaincode interface
func New() shim.Chaincode {
	return &scc{}
}

type scc struct{}

// Init implements the chaincode shim interface
func (s *scc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke implements the chaincode shim interface
func (s *scc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//fn, args := stub.GetFunctionAndParameters()
	fn, _ := stub.GetFunctionAndParameters()

	if fn == "getInputmaskIdx" {
		key := "inputmask"
		_value, _ := stub.GetState(key)
		if _value == nil {
			_value = []byte("0")
		}
		value, _ := strconv.Atoi(string(_value))
		stub.PutState(key, []byte(strconv.Itoa(value + 1)))
		return shim.Success([]byte(strconv.Itoa(value)))
	}
	//else if fn == "query" {
	//	if len(args) != 1 {
	//		return shim.Error("invalid number of args for query function")
	//	}
	//	key := args[0]
	//	value, err := stub.GetState(key)
	//	if err != nil {
	//		return shim.Error(fmt.Sprint("failed to get value for key %s", key))
	//	}
	//	return shim.Success([]byte(value))
	//} else if fn == "update" {
	//	if len(args) != 2 {
	//		return shim.Error("invalid number of args for update function")
	//	}
	//	key, value := args[0], []byte(args[1])
	//	err := stub.PutState(key, value)
	//	if err != nil {
	//		return shim.Error(fmt.Sprint("failed to update for key %s value %s", key, value))
	//	}
	//	return shim.Success([]byte("key-value updated"))
	//}
	return shim.Error("invalid function name.")
}

func main() {}
