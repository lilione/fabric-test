package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type eq struct{}

type eqtest struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	End        int    `json:"end"`
	M1         string `json:"m1"`
	M2         string `json:"m2"`
	U1         string `json:"u1"`
	U2         string `json:"u2"`
	Time       int    `json:"time"`
	Result     string `json:"result"`
}

// Init
func (t *eq) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func createCell(stub shim.ChaincodeStubInterface, key string, userID string, namespace string) {
	chainCodeArgs := ToChaincodeArgs("createCell", key, userID, namespace)
	response := stub.InvokeChaincode("honeybadgerscc", chainCodeArgs, "mychannel")

	if response.Status != shim.OK {
		fmt.Println(response.Message)
	}
}

func getActiveeqtests(stub shim.ChaincodeStubInterface) peer.Response {
	keysIter, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
	}
	defer keysIter.Close()

	var items []eqtest
	for keysIter.HasNext() {
		key, _ := keysIter.Next()
		item, _ := stub.GetState(key.Key)
		var eqtestInstance eqtest
		err = json.Unmarshal([]byte(item), &eqtestInstance)
		if eqtestInstance.Result == "None" {
			items = append(items, eqtestInstance)
		}
	}
	JSONasBytes, err := json.Marshal(items)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(JSONasBytes))

}

func getCompletedeqtests(stub shim.ChaincodeStubInterface) peer.Response {
	keysIter, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
	}
	defer keysIter.Close()

	var items []eqtest
	for keysIter.HasNext() {
		key, _ := keysIter.Next()
		item, _ := stub.GetState(key.Key)
		var eqtestInstance eqtest
		err = json.Unmarshal([]byte(item), &eqtestInstance)
		if eqtestInstance.Result != "None" {
			items = append(items, eqtestInstance)
		}
	}
	JSONasBytes, err := json.Marshal(items)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(JSONasBytes))

}

// Invoke
func (t *eq) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	if fn == "createeqtest" && len(args) >= 3 {
		fmt.Println("In createeqtest endpoint ")
		// Create a eqtest with the parameters provided
		// args[0] - name of the eqtest
		// args[1] - time to end the eqtest in seconds
		// args[2] - name of the user
		// returns memory cell id for the user move
		timeLimit, err := strconv.Atoi(args[1])
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType := "eqtest"

		memcell := "eq" + args[0] + args[2] + "cell" // in real use case this should
		// be a hash or something non-deterministic
		createCell(stub, memcell, args[2], "eq") // create a memory cell for the user to enter data
		eqtest := &eqtest{objectType, args[0], timeLimit, memcell, "None", args[2], "None", int(0), "None"}
		eqtestJSONasBytes, err := json.Marshal(eqtest)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState(args[0], eqtestJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(memcell))
	} else if fn == "getActiveeqtests" {

		return getActiveeqtests(stub)

	} else if fn == "getCompletedeqtests" {

		return getCompletedeqtests(stub)

	} else if fn == "joineqtest" && len(args) >= 2 {
		fmt.Println("In joineqtest endpoint")
		// Join a eqtest
		// args[0] - name of the eqtest
		// args[1] - name of the user
		// returns memory cell id for the user move
		memcell := "eq" + args[0] + args[1] + "cell"
		var eqtestInstance eqtest
		eqtestAsBytes, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal([]byte(eqtestAsBytes), &eqtestInstance)
		if err != nil {
			return shim.Error(err.Error())
		}
		createCell(stub, memcell, args[1], "eq")
		eqtestInstance.M2 = memcell
		eqtestInstance.U2 = args[1]
		eqtestJSONasBytes, err := json.Marshal(eqtestInstance)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState(args[0], []byte(eqtestJSONasBytes))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(memcell))
	} else if fn == "runMPC" && len(args) >= 1 {
		var eqtestInstance eqtest
		fmt.Println("In  run mpc operation")
		eqtestAsBytes, _ := stub.GetState(args[0])
		json.Unmarshal([]byte(eqtestAsBytes), &eqtestInstance)
		var cell1, cell2 secretcell
		cell1 = getCellMetaData(stub, eqtestInstance.M1)
		cell2 = getCellMetaData(stub, eqtestInstance.M2)
		if cell1.IsWritten && cell2.IsWritten {
			fmt.Println("Started MPC operation")
			mpcOp(stub, "equals", eqtestInstance.Name, eqtestInstance.M1, eqtestInstance.M2)
			return shim.Success([]byte("Started"))
		}
		return shim.Success([]byte("None"))
	} else if fn == "endeqtest" && len(args) >= 2 {
		// end the eqtest
		// args[0] - name of the eqtest
		var eqtestInstance eqtest
		fmt.Println("In endeqtest ")
		eqtestAsBytes, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal([]byte(eqtestAsBytes), &eqtestInstance)
		if err != nil {
			return shim.Error(err.Error())
		}
		// check if both cells are ready to be opened
		if eqtestInstance.M1 == "None" || eqtestInstance.M2 == "None" {
			return shim.Success([]byte("None"))
		}
		for {
			if checkMPCResult(stub, eqtestInstance.Name, "eq") != "None" {
				break
			}
			time.Sleep(2 * time.Second)
		}
		eqtestInstance.Result = checkMPCResult(stub, eqtestInstance.Name, "eq")
		eqtestJSONasBytes, err := json.Marshal(eqtestInstance)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState(args[0], []byte(eqtestJSONasBytes))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(eqtestInstance.Result))
	}
	return shim.Success([]byte("Invalid endpoint"))
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(eq)); err != nil {
		fmt.Printf("Error starting eq chaincode: %s", err)
	}
}
