package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	stateStartGlobal = "startGlobal"
	stateStartLocal = "startLocal"
	stateFinalizeGlobal = "finalizeGlobal"
	stateFinalizeLocal = "finalizeLocal"
)

type SmartContract struct {
}

type shipment struct {
	IdxInputProvider  	string
	IdxOutputProvider 	string
	//IdxAmount         	string
	Prev              	string
	Succs             	[]string
	State				string
}

type inquiry struct {
	Value 	string
	State	string
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

func getInquiry(stub shim.ChaincodeStubInterface, key string) inquiry {
	inquiryJSON, _ := stub.GetState(key)
	var inquiryInstance inquiry
	json.Unmarshal(inquiryJSON, &inquiryInstance)
	return inquiryInstance
}

func putInquiry(stub shim.ChaincodeStubInterface, key string, inquiryInstance inquiry) {
	inquiryJSON, _ := json.Marshal(&inquiryInstance)
	stub.PutState(key, inquiryJSON)
}
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "queryShipment" {
		itemID := args[0]
		seq := args[1]

		shipmentJSON, _ := stub.GetState(("itemInfo" + itemID + seq))

		return shim.Success(shipmentJSON)

	} else if fn == "queryInquiry" {
		itemID := args[0]
		seq := args[1]

		inquiryJSON, _ := stub.GetState("sourceItem" + itemID + seq)

		return shim.Success(inquiryJSON)

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

	} else if fn == "registerItemFinalizeGlobal" {
		idxRegistrant := args[0]
		//idxAmt := args[1]

		itemID := getCounter(stub,"itemID", 1)
		seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)

		var shipmentInstance shipment
		shipmentInstance.IdxOutputProvider = idxRegistrant
		//shipmentInstance.IdxAmount = idxAmt
		shipmentInstance.State = stateFinalizeGlobal
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		return shim.Success([]byte(itemID + " " + seq))

	} else if fn == "registerItemFinalizeLocal" {
		itemID := args[0]
		seq := args[1]
		maskedRegistrant := args[2]
		//maskedAmt := args[3]

		shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
		if shipmentInstance.State != stateFinalizeGlobal && shipmentInstance.State != stateFinalizeLocal {
			return shim.Error("registerItemGlobal not finished yet")
		}
		shipmentInstance.State = stateFinalizeLocal
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)
		idxRegistrant := shipmentInstance.IdxOutputProvider
		//idxAmt := shipmentInstance.IdxAmount

		shareRegistrant := calcShare(stub, idxRegistrant, maskedRegistrant)
		//shareAmt := calcShare(stub, idxAmt, maskedAmt)
		dbPut(stub, idxRegistrant, shareRegistrant)
		//dbPut(stub, idxAmt, shareAmt)

		return shim.Success([]byte(fmt.Sprintf("Register item %s finished", itemID)))

	} else if fn == "handOffItemStartGlobal" {
		itemID := args[0]
		prevSeq := args[1]

		prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))
		if prevShipmentInstance.State != stateFinalizeLocal {
			return shim.Error("Previous shipment not recorded yet")
		}

		seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)

		var shipmentInstance shipment
		shipmentInstance.State = stateStartGlobal
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		return shim.Success([]byte(seq))

	} else if fn == "handOffItemStartLocal" {
		idxInputProvider := args[0]
		maskedInputProvider := args[1]
		idxOutputProvider := args[2]
		maskedOutputProvider := args[3]
		itemID := args[4]
		prevSeq := args[5]
		seq := args[6]
		//idxAmt := args[4]
		//maskedAmt := args[5]
		//itemID := args[6]
		//prevSeq := args[7]
		//seq := args[8]

		prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))
		sharePrevOutputProvider := dbGet(stub, prevShipmentInstance.IdxOutputProvider)
		//sharePrevAmt := dbGet(stub, prevShipmentInstance.IdxAmount)

		shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
		if shipmentInstance.State != stateStartGlobal && shipmentInstance.State != stateStartLocal {
			return shim.Error("handOffItemStartGlobal not finished yet")
		}
		shipmentInstance.State = stateStartLocal
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		//handOffItem(stub, idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq, seq, sharePrevOutputProvider, sharePrevAmt)
		handOffItem(stub, idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, itemID, prevSeq, seq, sharePrevOutputProvider)

		return shim.Success([]byte(fmt.Sprintf("handOffItemStartLocal for item %s seq %s succeed", itemID, seq)))

	} else if fn == "handOffItemFinalizeGlobal" {
		idxInputProvider := args[0]
		idxOutputProvider := args[1]
		itemID := args[2]
		prevSeq := args[3]
		seq := args[4]
		//idxAmt := args[2]
		//itemID := args[3]
		//prevSeq := args[4]
		//seq := args[5]

		shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
		if shipmentInstance.State != stateStartLocal {
			return shim.Error("handOffItemStartLocal not finished yet")
		}
		shipmentInstance.IdxInputProvider = idxInputProvider
		shipmentInstance.IdxOutputProvider = idxOutputProvider
		//shipmentInstance.IdxAmount = idxAmt
		shipmentInstance.Prev = prevSeq
		shipmentInstance.State = stateFinalizeGlobal
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))
		prevShipmentInstance.Succs = append(prevShipmentInstance.Succs, seq)
		putShipment(stub, ("itemInfo" + itemID + prevSeq), prevShipmentInstance)

		return shim.Success([]byte(fmt.Sprintf("handOffItemFinalizeGlobal for item %s seq %s succeed", itemID, seq)))

	} else if fn == "handOffItemFinalizeLocal" {
		maskedInputProvider := args[0]
		maskedOutputProvider := args[1]
		itemID := args[2]
		seq := args[3]
		//maskedAmt := args[2]
		//itemID := args[3]
		//prevSeq := args[4]
		//seq := args[5]

		shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
		if shipmentInstance.State != stateFinalizeGlobal && shipmentInstance.State != stateFinalizeLocal {
			return shim.Error("handOffItemFinalizeGlobal not finished yet")
		}
		shipmentInstance.State = stateFinalizeLocal
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		idxInputProvider := shipmentInstance.IdxInputProvider
		idxOutputProvider := shipmentInstance.IdxOutputProvider
		//idxAmt := shipmentInstance.IdxAmount

		shareInputProvider := calcShare(stub, idxInputProvider, maskedInputProvider)
		shareOutputProvider := calcShare(stub, idxOutputProvider, maskedOutputProvider)
		//shareAmt := calcShare(stub, idxAmt, maskedAmt)
		dbPut(stub, idxInputProvider, shareInputProvider)
		dbPut(stub, idxOutputProvider, shareOutputProvider)
		//dbPut(stub, idxAmt, shareAmt)

		//prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))
		//sharePrevAmt := dbGet(stub, prevShipmentInstance.IdxAmount)
		//_sharePrevAmt, _ := strconv.Atoi(sharePrevAmt)
		//_shareAmt, _ := strconv.Atoi(shareAmt)
		//dbPut(stub, prevShipmentInstance.IdxAmount, strconv.Itoa(_sharePrevAmt - _shareAmt))

		return shim.Success([]byte(fmt.Sprintf("handOffItemFinalizeLocal for item %s seq %s succeed", itemID, seq)))

	} else if fn == "sourceItemStartLocal" {
		itemID := args[0]
		seq := args[1]

		shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
		if shipmentInstance.State != stateFinalizeLocal {
			return shim.Error("shipment not recorded yet")
		}

		shares := ""
		for true {
			shareInputProvider := dbGet(stub, shipmentInstance.IdxInputProvider)
			shares += shareInputProvider
			shipmentInstance = getShipment(stub, ("itemInfo" + itemID + shipmentInstance.Prev))
			if len(shipmentInstance.Prev) == 0 {
				break
			}
			shares += ","
		}

		var inquiryInstance inquiry
		inquiryInstance.State = stateStartLocal
		putInquiry(stub, ("sourceItem" + itemID + seq), inquiryInstance)

		sourceItem(stub, itemID, seq, shares)

		return shim.Success([]byte(fmt.Sprintf("sourceItemStartLocal for item %s seq %s succeed", itemID, seq)))
	} else if fn == "sourceItemFinalizeGlobal" {
		itemID := args[0]
		seq := args[1]
		listInputProvider := args[2]

		inquiryInstance := getInquiry(stub, ("sourceItem" + itemID + seq))
		if inquiryInstance.State != stateStartLocal {
			return shim.Error(fmt.Sprintf("sourceItemStartLocal for item %s seq %s not finished yet", itemID, seq))
		}

		inquiryInstance.Value = listInputProvider
		inquiryInstance.State = stateFinalizeGlobal
		putInquiry(stub, ("sourceItem" + itemID + seq), inquiryInstance)

		return shim.Success([]byte(fmt.Sprintf("sourceItemFinalizeGlobal for item %s seq %s succeed", itemID, seq)))
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
