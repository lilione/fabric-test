/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// New returns an implementation of the chaincode interface
func New() shim.Chaincode {
	return &scc{}
}

type scc struct {
}

type shipment struct {
	IdxLoadTime 	string
	IdxUnloadTime 	string
}

type truck struct {
	ShipmentList	[]shipment
}

// Init implements the chaincode shim interface
func (s *scc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func getCounter(stub shim.ChaincodeStubInterface, key string, inc int) string {
	cnt, _ := stub.GetState(key)
	if cnt == nil {
		cnt = []byte("0")
	}
	_cnt, _ := strconv.Atoi(string(cnt))
	stub.PutState(key, []byte(strconv.Itoa(_cnt + inc)))
	return string(cnt)
}

func getTruck(stub shim.ChaincodeStubInterface, key string) truck {
	truckJSON, _ := stub.GetState(key)
	var truckInstance truck
	json.Unmarshal(truckJSON, &truckInstance)
	return truckInstance
}

func putTruck(stub shim.ChaincodeStubInterface, key string, truckInstance truck) {
	truckJSON, _ := json.Marshal(&truckInstance)
	stub.PutState(key, truckJSON)
}

//func inRange(x string, l string, r string) bool {
//	if cmp(x, l) {
//		return false
//	}
//	if cmp(r, x) {
//		return false
//	}
//	return true
//}
//
//func intersect(l_1 string, r_1 string, l_2 string, r_2 string) bool {
//	var l_max string
//	if cmp(l_1, l_2) {
//		l_max = l_2
//	} else {
//		l_max = l_1
//	}
//	var r_min string
//	if cmp(r_1, r_2) {
//		r_min = r_1
//	} else {
//		r_min = r_2
//	}
//	if cmp(r_min, l_max) {
//		return false
//	}
//	return true
//}

// Invoke implements the chaincode shim interface
func (s *scc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "getInputmaskIdx" {
		num, _ := strconv.Atoi(args[0])

		var inputmaskIdx string
		value, _ := strconv.Atoi(getCounter(stub,"inputmaskCnt", num))
		for num > 0 {
			inputmaskIdx += strconv.Itoa(value)
			value += 1
			num -= 1
			if num > 0 {
				inputmaskIdx += " "
			}
		}

		return shim.Success([]byte(inputmaskIdx))

	} else if fn == "createTruck" {

		truckID := getCounter(stub, "truckID", 1)

		var truckInstance truck
		putTruck(stub, ("truckRegistry" + truckID), truckInstance)

		return shim.Success([]byte(truckID))

	} else if fn == "recordShipment" {

		truckID := args[0]
		idxLoadTime := args[1]
		maskedLoadTime := args[2]
		idxUnloadTime := args[3]
		maskedUnloadTime := args[4]

		shareLoadTime := calcShare(idxLoadTime, maskedLoadTime)
		shareUnloadTime := calcShare(idxUnloadTime, maskedUnloadTime)

		dbPut(idxLoadTime, shareLoadTime)
		dbPut(idxUnloadTime, shareUnloadTime)

		if !cmp(shareUnloadTime, shareLoadTime) {
			shipmentInstance := shipment{
				IdxLoadTime: idxLoadTime,
				IdxUnloadTime: idxUnloadTime,
			}

			truckInstance := getTruck(stub, ("truckRegistry" + truckID))
			truckInstance.ShipmentList = append(truckInstance.ShipmentList, shipmentInstance)
			putTruck(stub, ("truckRegistry" + truckID), truckInstance)

			return shim.Success([]byte("recorded successfully"))
		}

		return shim.Error("invalid load and unlaod time")

	} else if fn == "queryPositions" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]

		shareInitTime := calcShare(idxInitTime, maskedInitTime)
		shareEndTime := calcShare(idxEndTime, maskedEndTime)

		positions := ""
		truckInstance := getTruck(stub, ("truckRegistry" + truckID))
		for index, shipmentInstance := range truckInstance.ShipmentList {
			shareLoadTime := dbGet(shipmentInstance.IdxLoadTime)
			shareUnloadTime := dbGet(shipmentInstance.IdxUnloadTime)

			fmt.Println(index)
			if !cmp(shareInitTime, shareLoadTime) && !cmp(shareUnloadTime, shareEndTime) {
				if positions != "" {
					positions += " "
				}
				positions += strconv.Itoa(index)
			}
		}
		return shim.Success([]byte(positions))

	} else if fn == "queryNumber" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]

		shareInitTime := calcShare(idxInitTime, maskedInitTime)
		shareEndTime := calcShare(idxEndTime, maskedEndTime)

		cnt := "{0}"
		truckInstance := getTruck(stub, ("truckRegistry" + truckID))
		for index, shipmentInstance := range truckInstance.ShipmentList {
			fmt.Println("index", index)
			shareLoadTime := dbGet(shipmentInstance.IdxLoadTime)
			shareUnloadTime := dbGet(shipmentInstance.IdxUnloadTime)

			cnt = addShare(cnt, mulShare(oneMinusShare(cmpShare(shareInitTime, shareLoadTime)), oneMinusShare(cmpShare(shareUnloadTime, shareEndTime))))
			fmt.Println("cnt", cnt)
		}
		cnt = reconstruct(cnt)
		return shim.Success([]byte(cnt))
	}

	return shim.Error("invalid function name.")
}

func main() {
}
