package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	stateFinalized = "finalized"
	stateOngoing   = "ongoing"
)

type SmartContract struct {
}

type shipment struct {
	IdxInputProvider  	string
	IdxOutputProvider 	string
	IdxAmount         	string
	Prev              	string
	Succs             	[]string
	State				string
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
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

func getShipment(stub shim.ChaincodeStubInterface, key string) shipment {
	shipmentJSON, _ := stub.GetState(key)
	var shipmentInstance shipment
	json.Unmarshal(shipmentJSON, &shipmentInstance)
	return shipmentInstance
}

func putShipment(stub shim.ChaincodeStubInterface, key string, shipmentInstance shipment) {
	shipmentJSON, _ := json.Marshal(&shipmentInstance)
	stub.PutState(key, shipmentJSON)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "queryShipment" {
		itemID := args[0]
		seq := args[1]

		shipmentJSON, _ := stub.GetState(("itemInfo" + itemID + seq))

		return shim.Success(shipmentJSON)

	} else if fn == "getInputmaskIdx" {
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

	} else if fn == "registerItem" {
		idxRegistrant := args[0]
		maskedRegistrant := args[1]
		idxAmt := args[2]
		maskedAmt := args[3]

		itemID := getCounter(stub,"itemID", 1)
		seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)

		shareRegistrant := calcShare(stub, idxRegistrant, maskedRegistrant)
		shareAmt := calcShare(stub, idxAmt, maskedAmt)
		dbPut(stub, idxRegistrant, shareRegistrant)
		dbPut(stub, idxAmt, shareAmt)

		var shipmentInstance shipment
		shipmentInstance.IdxOutputProvider = idxRegistrant
		shipmentInstance.IdxAmount = idxAmt
		shipmentInstance.State = stateFinalized
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		return shim.Success([]byte(itemID + " " + seq))

	} else if fn == "handOffItemToNextProvider" {
		idxInputProvider := args[0]
		maskedInputProvider := args[1]
		idxOutputProvider := args[2]
		maskedOutputProvider := args[3]
		idxAmt := args[4]
		maskedAmt := args[5]
		itemID := args[6]
		prevSeq := args[7]

		prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))
		if prevShipmentInstance.State != stateFinalized {
			shim.Error("Previous shipment not finalized yet")
		}
		sharePrevOutputProvider := dbGet(stub, prevShipmentInstance.IdxOutputProvider)
		sharePrevAmt := dbGet(stub, prevShipmentInstance.IdxAmount)

		seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)

		var shipmentInstance shipment
		shipmentInstance.State = stateOngoing
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		handOffItemToNextProvider(stub, idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq, seq, sharePrevOutputProvider, sharePrevAmt)

		return shim.Success([]byte(seq))

	} else if fn == "handOffItemToNextProviderFinalize" {
		idxInputProvider := args[0]
		maskedInputProvider := args[1]
		idxOutputProvider := args[2]
		maskedOutputProvider := args[3]
		idxAmt := args[4]
		maskedAmt := args[5]
		itemID := args[6]
		prevSeq := args[7]
		seq := args[8]

		shareInputProvider := calcShare(stub, idxInputProvider, maskedInputProvider)
		shareOutputProvider := calcShare(stub, idxOutputProvider, maskedOutputProvider)
		shareAmt := calcShare(stub, idxAmt, maskedAmt)
		dbPut(stub, idxInputProvider, shareInputProvider)
		dbPut(stub, idxOutputProvider, shareOutputProvider)
		dbPut(stub, idxAmt, shareAmt)

		var shipmentInstance shipment
		shipmentInstance.IdxInputProvider = idxInputProvider
		shipmentInstance.IdxOutputProvider = idxOutputProvider
		shipmentInstance.IdxAmount = idxAmt
		shipmentInstance.Prev = prevSeq
		shipmentInstance.State = stateFinalized
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))
		sharePrevAmt := dbGet(stub, prevShipmentInstance.IdxAmount)
		_sharePrevAmt, _ := strconv.Atoi(sharePrevAmt)
		_shareAmt, _ := strconv.Atoi(shareAmt)
		dbPut(stub, prevShipmentInstance.IdxAmount, strconv.Itoa(_sharePrevAmt - _shareAmt))
		prevShipmentInstance.Succs = append(prevShipmentInstance.Succs, seq)
		putShipment(stub, ("itemInfo" + itemID + prevSeq), prevShipmentInstance)

		return shim.Success([]byte("Shipment details recorded"))

	}
	//else if fn == "sourceItem" {
	//	itemID := args[0]
	//	seq := args[1]
	//	response := sourceItem(stub, itemID, seq)
	//	if response.Status == 200 {
	//		return shim.Success(response.Payload)
	//	} else {
	//		return shim.Error(response.Message)
	//	}
	//}

	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
