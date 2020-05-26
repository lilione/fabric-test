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
	fn, args := stub.GetFunctionAndParameters()

	if fn == "getInputmaskIdx" {
		key := "inputmask"
		_value, _ := stub.GetState(key)
		if _value == nil {
			_value = []byte("0")
		}
		value, _ := strconv.Atoi(string(_value))

		stub.PutState(key, []byte(strconv.Itoa(value + 1)))

		return shim.Success([]byte(strconv.Itoa(value)))
	} else if fn == "sendMaskedInput" {
		idx := args[0]
		maskedInput := args[1]

		storeInput(idx, maskedInput)
		return shim.Success([]byte("success"))
	} else if fn == "reconstruct" {
		idx := args[0]

		result := reconstruct(idx)
		return shim.Success([]byte(result))
	}
	return shim.Error("invalid function name.")
}

func main() {}
