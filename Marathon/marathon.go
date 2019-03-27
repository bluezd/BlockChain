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

type ParticipantInfo struct {
	User_ID         string `json:"user_id"`
	User_Name       string `json:"user_name"`
	Birthday        string `json:"birthday"`
	National_ID     string `json:"national_id"`
	Passport_Number string `json:"passport_number"`
	Mobile          string `json:"mobile"`
	Point           string `json:"point_free"`
}

type MatchInfo struct {
	Match_ID   string `json:"match_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Match_Date string `json:"match_date"`
}

type MatchEnrollScoreInfo struct {
	User_Enter_ID string `json:"user_enter_id"`
	User_ID       string `json:"user_id"`
	Match_ID      string `json:"match_id"`
	Status        string `json:"status"`
	Match_Result  string `json:"match_result"`
	Score         string `json:"score"`
}

type MarathonChaincode struct {
}

// Init initializes chaincode
// ===========================
func (t *MarathonChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *MarathonChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("## invoke is running function: %s ##\n", function)

	// Handle different functions
	if function == "addParticipantInfo" {
		return addParticipantInfo(stub, args)
	} else if function == "updateParticipantInfo" {
		return updateParticipantInfo(stub, args)
	} else if function == "queryParticipantPoint" {
		return queryParticipantPoint(stub, args)
	} else if function == "addMatchEnrollScoreInfo" {
		return addMatchEnrollScoreInfo(stub, args)
	} else if function == "updateMatchEnrollScoreInfo" {
		return updateMatchEnrollScoreInfo(stub, args)
	} else if function == "queryParticipantInfo" || function == "queryMatchEnrollScoreInfo" ||
		function == "queryMatchInfo" {
		return queryHelper(function, stub, args)
	} else if function == "queryHistoryParticipantInfo" || function == "queryHistoryMatchEnrollScoreInfo" ||
		function == "queryHistoryMatchInfo" {
		return queryHistoryHelper(function, stub, args)
	} else if function == "queryMatchInfoBasedOnUser" {
		return queryMatchInfoBasedOnUser(stub, args)
	} else if function == "addMatchInfo" {
		return addMatchInfo(stub, args)
	} else if function == "updateMatchInfo" {
		return updateMatchInfo(stub, args)
	}

	fmt.Printf("!! invalid function: %s !!", function)
	return shim.Error("Received unknown function invocation")
}

func addMatchInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 4 {
		return shim.Error("!! Incorrect number of arguments, Expecting 4 !!")
	}

	// ==== Input Check ====
	fmt.Println("- start addMatchInfo -")

	// construct the key
	key := "MatchInfo_" + args[0]

	currentTime := timeHelper()
	fmt.Printf("[%s] <add> MatchInfo  %s", currentTime, key)

	ElementsAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get record: " + err.Error())
	} else if ElementsAsBytes != nil {
		return shim.Error("Match info record already exists: " + key)
	}

	matchInfoRecord := &MatchInfo{
		Match_ID:   args[0],
		Name:       args[1],
		Status:     args[2],
		Match_Date: args[3],
	}

	var matchInfoAsBytes []byte
	if matchInfoAsBytes, err = json.Marshal(matchInfoRecord); err != nil {
		return shim.Error(err.Error())
	}

	// Add Record to State
	if err = stub.PutState(key, matchInfoAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end addMatchInfo")
	return shim.Success(nil)
}

func addMatchEnrollScoreInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 6 {
		return shim.Error("!! Incorrect number of arguments, Expecting 6 !!")
	}

	// ==== Input Check ====
	fmt.Println("- start addMatchEnrollScoreInfo -")

	// construct the key
	key := "MatchEnrollScoreInfo_" + args[0]

	currentTime := timeHelper()
	fmt.Printf("[%s] <add> MatchEnrollScoreInfo  %s", currentTime, key)

	ElementsAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get record: " + err.Error())
	} else if ElementsAsBytes != nil {
		return shim.Error("Match record already exists: " + key)
	}

	matchEnrollScoreRecord := &MatchEnrollScoreInfo{
		User_Enter_ID: args[0],
		User_ID:       args[1],
		Match_ID:      args[2],
		Status:        args[3],
		Match_Result:  args[4],
		Score:         args[5],
	}

	var matchEnrollScoreInfoAsBytes []byte
	if matchEnrollScoreInfoAsBytes, err = json.Marshal(matchEnrollScoreRecord); err != nil {
		return shim.Error(err.Error())
	}

	// Add Record to State
	if err = stub.PutState(key, matchEnrollScoreInfoAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end addMatchEnrollScoreInfo")
	return shim.Success(nil)
}

func addParticipantInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var userID int

	if len(args) != 7 {
		return shim.Error("!! Incorrect number of arguments, Expecting 7 !!")
	}

	// ==== Input Check ====
	fmt.Println("- start addParticipantInfo -")
	if userID, err = strconv.Atoi(args[0]); err != nil {
		return shim.Error("1st argument user id must be a non-empty string")
	}

	currentTime := timeHelper()
	fmt.Printf("[%s] <add> Participant  %d", currentTime, userID)

	// construct the key
	key := "ParticipantInfo_" + args[0]

	// Check the Record in State
	ParticipantInfoAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get participant record: " + err.Error())
	} else if ParticipantInfoAsBytes != nil {
		return shim.Error("Participant record already exists: " + key)
	}

	participantRecord := &ParticipantInfo{
		User_ID:         args[0],
		User_Name:       args[1],
		Birthday:        args[2],
		National_ID:     args[3],
		Passport_Number: args[4],
		Mobile:          args[5],
		Point:           args[6],
	}

	var ParticipantInfoRecordAsBytes []byte
	if ParticipantInfoRecordAsBytes, err = json.Marshal(participantRecord); err != nil {
		return shim.Error(err.Error())
	}

	// Add Record to State
	if err = stub.PutState(key, ParticipantInfoRecordAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end addParticipantInfo")
	return shim.Success(nil)
}

func updateMatchInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input Check ====
	fmt.Println("- start updateMatchInfo -")
	if len(args) != 4 {
		return shim.Error("!! Incorrect number of arguments, Expecting 4 !!")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument user id must be a non-empty string")
	}

	// construct the key
	key := "MatchInfo_" + args[0]

	// Check the Record in State
	ElementsAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get match info record: " + err.Error())
	} else if ElementsAsBytes == nil {
		return shim.Error("match info record does not exists: " + key)
	}

	matchInfoRecord := &MatchInfo{}
	if err = json.Unmarshal(ElementsAsBytes, matchInfoRecord); err != nil {
		return shim.Error(err.Error())
	}

	matchInfoRecord.Name = args[1]
	matchInfoRecord.Status = args[2]
	matchInfoRecord.Match_Date = args[3]

	matchInfoRecordJSONBytes, err := json.Marshal(matchInfoRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	currentTime := timeHelper()
	fmt.Printf("[%s] <update>  updateMatchInfo %s", currentTime, key)

	// Add Record to State
	if err = stub.PutState(key, matchInfoRecordJSONBytes); err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- updateMatchInfo end ")
	return shim.Success(nil)
}

func updateMatchEnrollScoreInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input Check ====
	fmt.Println("- start updateMatchEnrollScoreInfo -")
	if len(args) != 6 {
		return shim.Error("!! Incorrect number of arguments, Expecting 6 !!")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument user id must be a non-empty string")
	}

	// construct the key
	key := "MatchEnrollScoreInfo_" + args[0]

	// Check the Record in State
	ElementsAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get match enroll record: " + err.Error())
	} else if ElementsAsBytes == nil {
		return shim.Error("match enroll record does not exists: " + key)
	}

	matchEnrollScoreRecord := &MatchEnrollScoreInfo{}
	if err = json.Unmarshal(ElementsAsBytes, matchEnrollScoreRecord); err != nil {
		return shim.Error(err.Error())
	}

	matchEnrollScoreRecord.User_ID = args[1]
	matchEnrollScoreRecord.Match_ID = args[2]
	matchEnrollScoreRecord.Status = args[3]
	matchEnrollScoreRecord.Match_Result = args[4]
	matchEnrollScoreRecord.Score = args[5]

	matchEnrollScoreRecordJSONBytes, err := json.Marshal(matchEnrollScoreRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	currentTime := timeHelper()
	fmt.Printf("[%s] <update> MatchEnrollScoreInfo %s", currentTime, key)

	// Add Record to State
	if err = stub.PutState(key, matchEnrollScoreRecordJSONBytes); err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateMatchEnrollScoreInfo")
	return shim.Success(nil)
}

func updateParticipantInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input Check ====
	fmt.Println("- start updateParticipantInfo -")
	if len(args) != 7 {
		return shim.Error("!! Incorrect number of arguments, Expecting 7 !!")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument user id must be a non-empty string")
	}

	// construct the key
	key := "ParticipantInfo_" + args[0]

	// Check the Record in State
	ParticipantInfoAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get participant record: " + err.Error())
	} else if ParticipantInfoAsBytes == nil {
		return shim.Error("Participant record does not exists: " + key)
	}

	participantRecord := &ParticipantInfo{}
	if err = json.Unmarshal(ParticipantInfoAsBytes, participantRecord); err != nil {
		return shim.Error(err.Error())
	}

	participantRecord.User_Name = args[1]
	participantRecord.Birthday = args[2]
	participantRecord.National_ID = args[3]
	participantRecord.Passport_Number = args[4]
	participantRecord.Mobile = args[5]
	participantRecord.Point = args[6]

	participantRecordJSONBytes, err := json.Marshal(participantRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	currentTime := timeHelper()
	fmt.Printf("[%s] <update> Participant %s", currentTime, key)

	// Add Record to State
	if err = stub.PutState(key, participantRecordJSONBytes); err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateParticipantInfo")
	return shim.Success(nil)
}

func getParticipantPointBytes(ParticipantInfoRecordAsBytes []byte) []byte {
	participantRecord := &ParticipantInfo{}
	_ = json.Unmarshal(ParticipantInfoRecordAsBytes, participantRecord)

	participantPointRecord := struct {
		User_ID string `json:"user_id"`
		Point   string `json:"point"`
	}{
		User_ID: participantRecord.User_ID,
		Point:   participantRecord.Point,
	}

	participantPointRecordJSONBytes, _ := json.Marshal(participantPointRecord)

	return participantPointRecordJSONBytes
}

func queryHelper(funcName string, stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Printf("- start queryHelper: %s\n", args[0])
	if len(args) != 1 {
		return shim.Error("!! Incorrect number of arguments, Expecting 1 !!")
	}

	// construct the key
	key := strings.TrimLeft(funcName, "query") + "_" + args[0]

	queryInfoAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		fmt.Println(jsonResp)
		return shim.Error(jsonResp)
	} else if queryInfoAsBytes == nil {
		jsonResp := "{\"Error\":\" ID does not exist: " + key + "\"}"
		fmt.Println(jsonResp)
		return shim.Error(jsonResp)
	}

	fmt.Println("- end queryHelper")
	return shim.Success(queryInfoAsBytes)
}

func queryParticipantPoint(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Printf("- start queryParticipantPoint: %s\n", args[0])
	if len(args) != 1 {
		return shim.Error("!! Incorrect number of arguments, Expecting 1 !!")
	}

	// construct the key
	key := "ParticipantInfo_" + args[0]

	ParticipantInfoRecordAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		fmt.Println(jsonResp)
		return shim.Error(jsonResp)
	} else if ParticipantInfoRecordAsBytes == nil {
		jsonResp := "{\"Error\":\" ID does not exist: " + key + "\"}"
		fmt.Println(jsonResp)
		return shim.Error(jsonResp)
	}

	ParticipantInfoRecordAsBytes = getParticipantPointBytes(ParticipantInfoRecordAsBytes)

	fmt.Println("- end queryParticipantPoint")
	return shim.Success(ParticipantInfoRecordAsBytes)
}

func queryHistoryHelper(funcName string, stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("!! Incorrect number of arguments, Expecting 1 !!")
	}

	// construct the key
	key := strings.TrimLeft(funcName, "queryHistory") + "_" + args[0]

	fmt.Printf("- start queryHistoryHelper: %s\n", key)

	resultsIterator, err := stub.GetHistoryForKey(key)
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

	fmt.Printf("- end queryHistoryHelper:\n   %s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ===== Example: Parameterized rich query =================================================
// queryVehiclePartByOwner queries for vehicle part based on a passed in owner.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// =========================================================================================
func queryMatchInfoBasedOnUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userID := args[0]

	//This is for native fabric: couchdb
	//queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"vehiclePart\",\"owner\":\"%s\"}}", owner)

	//This is for bcs: bdb
	//queryString := fmt.Sprintf("SELECT valueJson FROM <STATE> WHERE json_extract(valueJson, '$.user_id', '$.score') = '[\"%s\",\"90\"]'", userID)
	queryString := fmt.Sprintf("SELECT valueJson FROM <STATE> WHERE json_extract(valueJson, '$.user_id') = '%s' AND json_extract(valueJson, '$.score') IS NOT NULL", userID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
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
	err := shim.Start(new(MarathonChaincode))
	if err != nil {
		fmt.Printf("Error starting Parts Trace chaincode: %s", err)
	}
}
