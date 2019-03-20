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
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type ParticipantInfo struct {
	User_Id         string `json:"user_id"`
	User_Name       string `json:"user_name"`
	Birthday        string `json:"birthday"`
	National_Id     string `json:"national_id"`
	Passport_Number string `json:"passport_number"`
	Mobile          string `json:"mobile"`
	Point           string `json:"point"`
}

type MatchInfo struct {
	Match_ID   string `json:"match_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Match_Date string `json:"match_date"`
}

type MatchEnrollInfo struct {
	User_Enter_Id string `json:"user_enter_id"`
	User_ID       string `json:"user_id"`
	Match_ID      string `json:"match_id"`
	Status        string `json:"status"`
	Match_Result  string `json:"match_result"`
	Score         string `json:"score"`
}

/*
type MatchScore struct {
	User_Enter_Id string `json:"user_enter_id"`
	User_ID       string `json:"user_id"`
	Match_ID      string `json:"match_id"`
	Match_Result  string `json:"match_result"`
	Score         string `json:"score"`
}
*/

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
	} else if function == "updateParticipantInfo" || function == "updateParticipantPoint" {
		return updateParticipantInfo(stub, args)
	} else if function == "queryParticipantInfo" {
		return queryParticipantInfo(stub, args, false)
	} else if function == "queryParticipantPoint" {
		return queryParticipantInfo(stub, args, true)
	} else if function == "queryHistoryParticipantInfo" {
		return queryHistoryParticipantInfo(stub, args, false)
	} else if function == "queryHistoryParticipantPoint" {
		return queryHistoryParticipantInfo(stub, args, true)
	}

	fmt.Printf("!! invalid function: %s !!", function)
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// CreateCertificate
// ============================================================
func addParticipantInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var userID int

	if len(args) != 6 {
		return shim.Error("!! Incorrect number of arguments, Expecting 6 !!")
	}

	// ==== Input Check ====
	fmt.Println("- start addParticipantInfo -")
	if userID, err = strconv.Atoi(args[0]); err != nil {
		return shim.Error("1st argument user id must be a non-empty string")
	}

	currentTime := timeHelper()
	fmt.Printf("[%s] <add> Participant  %d", currentTime, User_Id)

	// construct the key
	key := args[0]

	// Check the Record in State
	ParticipantInfoAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get certificate record: " + err.Error())
	} else if ParticipantInfoAsBytes != nil {
		return shim.Error("Participant record already exists: " + key)
	}

	participantRecord := &ParticipantInfo{
		User_ID:         args[0],
		User_Name:       args[1],
		Birthday:        args[2],
		National_Id:     args[3],
		Passport_Number: args[4],
		Mobile:          args[5],
	}
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

func updateParticipantInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input Check ====
	fmt.Println("- start updateParticipantInfo -")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument user id must be a non-empty string")
	}

	// construct the key
	key := args[0]

	// Check the Record in State
	ParticipantInfoAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get participant record: " + err.Error())
	} else if certificateRecordAsBytes == nil {
		return shim.Error("Participant record does not exists: " + key)
	}

	participantRecord := &ParticipantInfo{}
	if err = json.Unmarshal(ParticipantInfoAsBytes, participantRecord); err != nil {
		return shim.Error(err.Error())
	}

	if len(args) == 6 {
		participantRecord.User_Name = args[1]
		participantRecord.Birthday = args[2]
		participantRecord.National_Id = args[3]
		participantRecord.Passport_Number = args[4]
		participantRecord.Mobile = args[5]
	} else if len(args) == 2 {
		participantRecord.Point = args[6]
	} else {
		return shim.Error("!! Incorrect number of arguments, Expecting 6 !!")
	}

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
	if err = json.Unmarshal(ParticipantInfoRecordAsBytes, participantRecord); err != nil {
		return shim.Error(err.Error())
	}

	participantPointRecord := struct {
		User_Id string `json:"user_id"`
		Point   string `json:"point"`
	}{
		User_Id: participantRecord.User_Id,
		Point:   participantRecord.Point,
	}

	participantPointRecordJSONBytes, err := json.Marshal(participantPointRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	return participantPointRecordJSONBytes
}

func queryParticipantInfo(stub shim.ChaincodeStubInterface, args []string, pointQuery bool) []byte {
	var err error

	fmt.Printf("- start queryParticipantInfo: %s\n", args[0])
	if len(args) != 1 {
		return shim.Error("!! Incorrect number of arguments, Expecting 1 !!")
	}

	// construct the key
	key := args[0]

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

	if pointQuery {
		ParticipantInfoRecordAsBytes = getParticipantPointBytes(ParticipantInfoRecordAsBytes)
	}

	fmt.Println("- end queryParticipantInfo")
	return shim.Success(ParticipantInfoRecordAsBytes)
}

func queryHistoryParticipantInfo(stub shim.ChaincodeStubInterface, args []string, pointQuery bool) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("!! Incorrect number of arguments, Expecting 2 !!")
	}

	// construct the key
	key := args[0]

	fmt.Printf("- start queryHistoryParticipantInfo: %s\n", key)

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
			if pointQuery {
				buffer.WriteString(string(getParticipantPointBytes(response.Value)))

			} else {
				buffer.WriteString(string(response.Value))
			}
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

	fmt.Printf("- end queryHistoryParticipantInfo:\n   %s\n", buffer.String())

	return shim.Success(buffer.Bytes())
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
