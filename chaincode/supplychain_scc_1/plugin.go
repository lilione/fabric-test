/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric-protos-go/peer"
)

// New returns an implementation of the chaincode interface
func New() shim.Chaincode {
	return &scc{}
}

type scc struct {
}

// Init implements the chaincode shim interface
func (s *scc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke implements the chaincode shim interface
func (s *scc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	fmt.Println("scc", fn, args)

	if fn == "dbPut" {
		key := args[0]
		value := args[1]

		dbPut(key, value)

		return shim.Success([]byte(""))

	} else if fn == "dbGet" {
		key := args[0]

		value := dbGet(key)

		return shim.Success(value)

	} else if fn == "calcShare" {
		idx := args[0]
		maskedShare := args[1]

		value := calcShare(idx, maskedShare)

		return shim.Success([]byte(value))

	} else if fn == "handOffItem" {
		idxInputProvider := args[0]
		maskedInputProvider := args[1]
		idxOutputProvider := args[2]
		maskedOutputProvider := args[3]
		itemID := args[4]
		prevSeq := args[5]
		seq := args[6]
		sharePrevOutputProvider := args[7]
		//idxAmt := args[4]
		//maskedAmt := args[5]
		//itemID := args[6]
		//prevSeq := args[7]
		//seq := args[8]
		//sharePrevOutputProvider := args[9]
		//sharePrevAmt := args[10]

		handOffItem(idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, itemID, prevSeq, seq, sharePrevOutputProvider)

	} else if fn == "sourceItem" {
		itemID := args[0]
		seq := args[1]
		shares := args[2]

		sourceItem(itemID, seq, shares)

	}

	return shim.Error("invalid function name.")
}

func main() {}
