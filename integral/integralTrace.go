/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type IntegralChaincode struct {
}

type Integral struct {
	UserName       string `json:"userNme"`
	EnterpriseName string `json:"enterpriseName"`
	IntegralCount  int    `json:"integralCount"`
	AddNote        string `json:"addNote"`
}

// Init initializes chaincode
// ===========================
func (t *IntegralChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *IntegralChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("## invoke is running function: %s ##\n", function)

	// Handle different functions
	if function == "initIntegral" {
		return t.initIntegral(stub, args)
	} else if function == "addIntegral" {
		return t.addIntegral(stub, args)
	} else if function == "convertIntegral" {
		return t.convertIntegral(stub, args)
	} else if function == "queryIntegral" {
		return t.queryIntegral(stub, args)
	} else if function == "queryHistoryIntegral" {
		return t.queryHistoryIntegral(stub, args)
		//} else if function == "deleteIntegral" {
		//		return t.deleteIntegral(stub, args)
	}

	fmt.Printf("!! invalid function: %s !!", function)
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initVehiclePart - create a new vehicle part, store into chaincode state
// ============================================================
func (t *IntegralChaincode) initIntegral(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 4 {
		return shim.Error("!! Incorrect number of arguments, Expecting 4 !!")
	}

	// ==== Input Check ====
	fmt.Println("- start init integral -")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument username must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument enterprisename must be a non-empty string")
	}

	userName := strings.ToLower(args[0])
	enterpriseName := strings.ToLower(args[1])
	integralCount, _ := strconv.Atoi(args[2])
	addNote := args[3]

	currentTime := timeHelper()
	if len(args[3]) <= 0 {
		addNote = fmt.Sprintf("[%s] <init> %d integral", currentTime, integralCount)
	}

	if res := enterpriseNameCheck(enterpriseName); res != true {
		s := fmt.Sprintf("!! Invalid Enterprise Name: %s !!\n", enterpriseName)
		fmt.Printf(s)
		return shim.Error(s)
	}

	// construct the key
	keyComposite := userName + "-" + enterpriseName

	// Check the Record in State
	integralRecordAsBytes, err := stub.GetState(keyComposite)
	if err != nil {
		return shim.Error("Failed to get integral record: " + err.Error())
	} else if integralRecordAsBytes != nil {
		fmt.Println("integral record already exists: " + keyComposite)
		//return shim.Error("Integral UserName already exists: " + userName)
	}

	integralRecord := &Integral{userName, enterpriseName, integralCount, addNote}
	integralRecordAsBytes, err = json.Marshal(integralRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Add Record to State
	err = stub.PutState(keyComposite, integralRecordAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init integral")
	return shim.Success(nil)
}

func (t *IntegralChaincode) addIntegral(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 3 {
		return shim.Error("!! Incorrect number of arguments, Expecting 3 !!")
	}

	// ==== Input Check ====
	fmt.Println("- start add integral -")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument username must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument enterprisename must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3nd argument enterprisename must be a non-empty string")
	}

	userName := strings.ToLower(args[0])
	enterpriseName := strings.ToLower(args[1])
	integralCount, _ := strconv.Atoi(args[2])

	currentTime := timeHelper()

	if res := enterpriseNameCheck(enterpriseName); res != true {
		s := fmt.Sprintf("!! Invalid Enterprise Name: %s !!\n", enterpriseName)
		fmt.Printf(s)
		return shim.Error(s)
	}

	// construct the key
	keyComposite := userName + "-" + enterpriseName

	// Check the Record in State
	integralRecordAsBytes, err := stub.GetState(keyComposite)
	if err != nil {
		return shim.Error("Failed to get integral record: " + err.Error())
	} else if integralRecordAsBytes == nil {
		fmt.Printf("!! integral record does not exist: %s !!", keyComposite)
		return shim.Error("Integral UserName does not exist: " + keyComposite)
	}
	integralRecord := &Integral{}
	err = json.Unmarshal(integralRecordAsBytes, integralRecord)
	integralRecord.IntegralCount += integralCount
	integralRecord.AddNote = fmt.Sprintf("[%s] <add> %d integral", currentTime, integralCount)

	integralRecordJSONBytes, err := json.Marshal(integralRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Add Record to State
	err = stub.PutState(keyComposite, integralRecordJSONBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end add integral")
	return shim.Success(nil)
}

func (t *IntegralChaincode) queryIntegral(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("!! Incorrect number of arguments, Expecting 2 !!")
	}
	userName := strings.ToLower(args[0])
	enterpriseName := strings.ToLower(args[1])

	// construct the key
	keyComposite := userName + "-" + enterpriseName

	fmt.Printf("- start queryIntegral: %s\n", keyComposite)
	integralRecordAsBytes, err := stub.GetState(keyComposite)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + keyComposite + "\"}"
		fmt.Println(jsonResp)
		return shim.Error(jsonResp)
	} else if integralRecordAsBytes == nil {
		jsonResp := "{\"Error\":\" user enterprise does not exist: " + keyComposite + "\"}"
		fmt.Println(jsonResp)
		return shim.Error(jsonResp)
	}

	fmt.Println("- end queryIntegral")
	return shim.Success(integralRecordAsBytes)
}

func (t *IntegralChaincode) queryHistoryIntegral(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("!! Incorrect number of arguments, Expecting 2 !!")
	}
	userName := strings.ToLower(args[0])
	enterpriseName := strings.ToLower(args[1])

	// construct the key
	keyComposite := userName + "-" + enterpriseName

	fmt.Printf("- start queryHistoryIntegral: %s\n", keyComposite)

	resultsIterator, err := stub.GetHistoryForKey(keyComposite)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON Integral)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("-  queryHistoryIntegral returning:\n   %s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (t *IntegralChaincode) convertIntegral(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 4 {
		return shim.Error("!! Incorrect number of arguments, Expecting 4 !!")
	}
	userName := strings.ToLower(args[0])
	origEnterpriseName := strings.ToLower(args[1])
	targetEnterpriseName := strings.ToLower(args[2])
	convertIntegralCount, _ := strconv.Atoi(args[3])

	if !enterpriseNameCheck(origEnterpriseName) || !enterpriseNameCheck(targetEnterpriseName) {
		s := fmt.Sprintf("!! Invalid Enterprise Name !!\n")
		fmt.Printf(s)
		return shim.Error(s)
	}

	// construct the key
	keyComposite := userName + "-" + origEnterpriseName

	fmt.Println("- start convert integral")
	// Get Orig Integral
	origIntegralRecordAsBytes, err := stub.GetState(keyComposite)
	if err != nil {
		return shim.Error("Failed to get integral record: " + err.Error())
	} else if origIntegralRecordAsBytes == nil {
		fmt.Printf("!! integral record does not exist: %s !!\n", keyComposite)
		return shim.Error("Integral UserName does not exist: " + keyComposite)
	}
	origIntegralRecord := new(Integral)
	err = json.Unmarshal(origIntegralRecordAsBytes, origIntegralRecord)
	if err != nil {
		fmt.Println("!! failed to unmarshall to origIntegralRecord !!")
		return shim.Error("Failed to unmarshall to origIntegralRecord")
	}

	// integral convert rate bank:telecom:shopping_mall = 1:2:4
	integralRateMap := map[string]int{
		"bank":          0,
		"telecom":       0,
		"shopping_mall": 0,
	}

	switch origIntegralRecord.EnterpriseName {
	case "bank":
		integralRateMap["bank"] = origIntegralRecord.IntegralCount - convertIntegralCount
		integralRateMap["telecom"] = convertIntegralCount * 2
		integralRateMap["shopping_mall"] = convertIntegralCount * 4
	case "telecom":
		integralRateMap["bank"] = convertIntegralCount / 2
		integralRateMap["telecom"] = origIntegralRecord.IntegralCount - convertIntegralCount
		integralRateMap["shopping_mall"] = convertIntegralCount * 2
		fmt.Println("   - telecom")
	case "shopping_mall":
		integralRateMap["bank"] = convertIntegralCount / 4
		integralRateMap["telecom"] = convertIntegralCount / 2
		integralRateMap["shopping_mall"] = origIntegralRecord.IntegralCount - convertIntegralCount
	default:
		fmt.Println("!! Incorrect Enterprise Name [bank, telecom, shopping_mall] !!")
		return shim.Error("Incorrect Enterprise Name [bank, telecom, shopping_mall]")
	}
	fmt.Printf("   - bank=%d, telecom=%d, shopping_mall=%d\n", integralRateMap["bank"], integralRateMap["telecom"], integralRateMap["shopping_mall"])

	currentTime := timeHelper()
	fmt.Printf("   - orig username=%s, enterprisename=%s, ingegral=%d\n", origIntegralRecord.UserName, origIntegralRecord.EnterpriseName, origIntegralRecord.IntegralCount)

	// update the orig State
	fmt.Printf("   - update the orig (%s) state\n", keyComposite)
	origIntegralRecord.IntegralCount = integralRateMap[origIntegralRecord.EnterpriseName]
	origIntegralRecord.AddNote = fmt.Sprintf("[%s] <convert> reduce %d integral", currentTime, convertIntegralCount)
	origIntegralRecordJSONBytes, err := json.Marshal(origIntegralRecord)
	if err != nil {
		fmt.Println("!! failed to unmarshall to origIntegralRecordJSONBytes !!")
		return shim.Error(err.Error())
	}
	fmt.Printf("   - origIntegralRecordJSONBytes: %s\n", origIntegralRecordJSONBytes)
	err = stub.PutState(keyComposite, origIntegralRecordJSONBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	//delete(integralRateMap, origIntegralRecord.EnterpriseName)

	// handle the target
	// construct the key
	keyComposite = userName + "-" + targetEnterpriseName

	// Get target Integral
	targetIntegralRecordAsBytes, err := stub.GetState(keyComposite)
	if err != nil {
		return shim.Error("Failed to get integral record: " + err.Error())
	} else if targetIntegralRecordAsBytes == nil {
		fmt.Printf("!! integral record does not exist: %s !!\n", keyComposite)
	}

	var targetIntegralRecord *Integral

	if targetIntegralRecordAsBytes != nil {
		targetIntegralRecord = new(Integral)
		err = json.Unmarshal(targetIntegralRecordAsBytes, targetIntegralRecord)
		targetIntegralRecord.IntegralCount += integralRateMap[targetEnterpriseName]
	} else {
		targetIntegralRecord = &Integral{userName, targetEnterpriseName, integralRateMap[targetEnterpriseName], ""}
	}
	targetIntegralRecord.AddNote = fmt.Sprintf("[%s]<convert> add %d integral\n", currentTime, integralRateMap[targetEnterpriseName])
	fmt.Printf("   - target username=%s, enterprisename=%s, ingegral=%d\n", targetIntegralRecord.UserName, targetIntegralRecord.EnterpriseName, targetIntegralRecord.IntegralCount)

	// update target Integral
	fmt.Printf("   - update the target (%s) state\n", keyComposite)
	targetIntegralRecordJSONBytes, err := json.Marshal(targetIntegralRecord)
	if err != nil {
		fmt.Println("!! failed to marshall to targetIntegralRecordJSONBytes !!")
		return shim.Error(err.Error())
	}
	err = stub.PutState(keyComposite, targetIntegralRecordJSONBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end convert integral")
	return shim.Success(nil)
}

/*
func (t *IntegralChaincode) deleteIntegral(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("!! Incorrect number of arguments, Expecting 2 !!")
	}

	fmt.Println("- start delete integral -")

	userName := strings.ToLower(args[0])
	enterpriseName := strings.ToLower(args[1])

	// construct the key
	keyComposite := userName + "-" + enterpriseName

	// Check the Record in state
	integralRecordAsBytes, err := stub.GetState(keyComposite)
	if err != nil {
		return shim.Error("Failed to get integral record: " + err.Error())
	} else if integralRecordAsBytes == nil {
		fmt.Println("integral record does not exist: " + keyComposite)
		return shim.Error("Integral UserName does not exist: " + keyComposite)
	}

	// remove the user from state
	err = stub.DelState(keyComposite)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to delete %s: %s\n", keyComposite, err.Error()))
	}

	fmt.Println("- end delete integral -")
	return shim.Success(nil)
}
*/

var enterpriseNameList = []string{"bank", "telecom", "shopping_mall"}

func enterpriseNameCheck(name string) bool {
	for _, value := range enterpriseNameList {
		if name == value {
			return true
		}
	}
	return false
}

func timeHelper() string {
	t := time.Now()
	currentTime := t.Format("2006-01-02T15:04:05")
	return currentTime
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(IntegralChaincode))
	if err != nil {
		fmt.Printf("Error starting Parts Trace chaincode: %s", err)
	}
}
