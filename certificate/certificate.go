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

type SmartContract struct {
}
type Certificate struct {
	CertificateHash   string `json:CertificateHash`
	PartnerName       string `json:PartnerName`
	Contacts          string `json:Contacts`
	Mobile            string `json:Mobile`
	Email             string `json:Email`
	CertificateType   string `json:CertificateType`
	CertificateName   string `json:CertificateName`
	PassingDate       string `json:PassingDate`
	ExpiryDate        string `json:ExpiryDate`
	CertificateStatus string `json:CertificateStatus` // 0:通过 1:失败 2:降级通过 3:取消
	Participant       string `json: Participant`
	Score             string `json: Score`
}

var CerfificationQueryMap = map[string]string{
	"PartnerName":     "partnername~all",
	"CertificateName": "certificate~all",
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting transaction Trace chaincode: %s", err)
	}
}
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)

}
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// Retrieve the requested Smart Contract function and arguments

	function, args := stub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately

	if function == "createCertificate" {
		return s.createCertificate(stub, args)
	} else if function == "getHistoryForRecord" {
		return s.getHistoryForRecord(stub, args)
	} else if function == "queryCertificate" {
		return s.queryCertificate(stub, args)
	} else if function == "removeCertificate" {
		return s.removeCertificate(stub, args)
	} else if function == "updateCertificate" {
		return s.updateCertificate(stub, args)
	} else if function == "queryCertificateBasedOnName" {
		return s.queryCertificateBasedOnName(stub, args)
	} else if function == "queryAllCertificate" {
		return s.queryAllCertificate(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name." + function)

}
func (s *SmartContract) queryCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")

	}
	key := args[0]

	certificateAsBytes, _ := stub.GetState(key)
	if certificateAsBytes == nil {
		return shim.Error("Could not locate Certificate" + key)

	}
	certificate := Certificate{}

	json.Unmarshal(certificateAsBytes, &certificate)

	// Normally check that the specified argument is a valid holder of tuna but here we are skipping this check for this example.

	certificateAsBytes, _ = json.Marshal(certificate)

	// err := stub.PutState(args[0], certificateAsBytes)

	// if err != nil {

	// 	return shim.Error(fmt.Sprintf("Failed to record accessor of certificate: %s", args[0]))

	// }

	return shim.Success(certificateAsBytes)
}
func (s *SmartContract) createCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12")
	}

	key := args[0]

	keyAsBytes, _ := stub.GetState(key)
	if keyAsBytes != nil {
		return shim.Error("Certificate key already exists:" + key)

	}

	var certificate = Certificate{CertificateHash: args[0], PartnerName: args[1], Contacts: args[2], Mobile: args[3], Email: args[4], CertificateType: args[5], CertificateName: args[6], PassingDate: args[7], ExpiryDate: args[8], CertificateStatus: args[9], Participant: args[10], Score: args[11]}

	certificateAsBytes, _ := json.Marshal(certificate)

	err := stub.PutState(args[0], certificateAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record certificate catch: %s", args[0]))
	}

	// create index
	err = createIndexHelper(stub, &certificate)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) removeCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	certificateAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get Certificate record: " + err.Error())
	} else if certificateAsBytes == nil {
		fmt.Println("Certificate does not exist: " + args[0])
		return shim.Error("Certificate does not exist: " + args[0])
	}

	certificate := &Certificate{}
	json.Unmarshal(certificateAsBytes, certificate)

	err = stub.DelState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to remove Certificate catch: %s", args[0]))
	}

	// delete index
	err = deleteIndexHelper(stub, certificate)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) updateCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12")
	}
	certificateAsBytes, _ := stub.GetState(args[0])

	if certificateAsBytes == nil {
		return shim.Error("Could not locate certificateAsBytes by" + args[0])
	}

	certificate := Certificate{}
	json.Unmarshal(certificateAsBytes, &certificate)
	// delete index
	err := deleteIndexHelper(stub, &certificate)
	if err != nil {
		return shim.Error(err.Error())
	}

	certificate.PartnerName = args[1]
	certificate.Contacts = args[2]
	certificate.Mobile = args[3]
	certificate.Email = args[4]
	certificate.CertificateType = args[5]
	certificate.CertificateName = args[6]
	certificate.PassingDate = args[7]
	certificate.ExpiryDate = args[8]
	certificate.CertificateStatus = args[9]
	certificate.Participant = args[10]
	certificate.Score = args[11]

	certificateAsBytes, _ = json.Marshal(certificate)
	err = stub.PutState(args[0], certificateAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to update certificate: %s", args[0]))
	}

	// create index
	err = createIndexHelper(stub, &certificate)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SmartContract) getHistoryForRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	recordKey := args[0]
	fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

	resultsIterator, err := stub.GetHistoryForKey(recordKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the key/value pair
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
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
		//as-is (as the Value itself a JSON vehiclePart)
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

	fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryAllCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	startKey := ""
	endKey := ""

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	fmt.Printf("- queryAllCertificates:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryCertificateBasedOnName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	queryName := args[1]
	fmt.Println("- start  queryCertificateBasedOnName", queryName)
	_, ok := CerfificationQueryMap[args[0]]
	if !ok {
		fmt.Println("!! Incorrect Query Option [PartnerName, CertificateName] !!")
		return shim.Error("Incorrect Query Option [PartnerName, CertificateName]")
	}

	certificateResultsIterator, err := stub.GetStateByPartialCompositeKey(CerfificationQueryMap[args[0]], []string{queryName})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer certificateResultsIterator.Close()

	var buffer bytes.Buffer
	var buffer1 bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	var i int
	for i = 0; certificateResultsIterator.HasNext(); i++ {
		buffer1.Reset()
		responseRange, err := certificateResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		var returnedPartnerName, returnedCertificateName string
		if args[0] == "PartnerName" {
			returnedPartnerName = compositeKeyParts[0]
			//returnedCertificateHash := compositeKeyParts[2]
			//returnedContacts := compositeKeyParts[3]
			//returnedMobile := compositeKeyParts[4]
			//returnedEmail := compositeKeyParts[5]
			//returnedCertificateType := compositeKeyParts[6]
			returnedCertificateName = compositeKeyParts[1]
			//returnedPassingDate := compositeKeyParts[7]
			//returnedExpiryDate := compositeKeyParts[8]
			//returnedCertificateStatus := compositeKeyParts[9]
			//returnedParticipant := compositeKeyParts[10]
			//returnedScore := compositeKeyParts[11]

			fmt.Printf("- found a certificate record from index:%s partnerName:%s certificateName:%s certificateHash:%s contacts:%s mobile:%s email:%s certificateType:%s passingDate:%s expiryDate:%s certificateStatus:%s participant:%s score:%s\n",
				objectType, compositeKeyParts[0], compositeKeyParts[1],
				compositeKeyParts[2], compositeKeyParts[3], compositeKeyParts[4], compositeKeyParts[5],
				compositeKeyParts[6], compositeKeyParts[7], compositeKeyParts[8],
				compositeKeyParts[9], compositeKeyParts[10], compositeKeyParts[11])
		} else if args[0] == "CertificateName" {
			returnedPartnerName = compositeKeyParts[1]
			returnedCertificateName = compositeKeyParts[0]

			fmt.Printf("- found a certificate record from index:%s certificateName:%s partnerName:%s certificateHash:%s contacts:%s mobile:%s email:%s certificateType:%s passingDate:%s expiryDate:%s certificateStatus:%s participant:%s score:%s\n",
				objectType, compositeKeyParts[0], compositeKeyParts[1],
				compositeKeyParts[2], compositeKeyParts[3], compositeKeyParts[4], compositeKeyParts[5],
				compositeKeyParts[6], compositeKeyParts[7], compositeKeyParts[8],
				compositeKeyParts[9], compositeKeyParts[10], compositeKeyParts[11])
		}

		buffer1.WriteString("{\"CertificateHash\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[2])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"PartnerName\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(returnedPartnerName)
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"Contacts\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[3])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"Mobile\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[4])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"Email\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[5])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"CertificateType\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[6])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"CertificateName\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(returnedCertificateName)
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"PassingDate\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[7])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"ExpiryDate\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[8])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"CertificateStatus\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[9])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"Participant\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[10])
		buffer1.WriteString("\"")

		buffer1.WriteString(", \"Score\":")
		buffer1.WriteString("\"")
		buffer1.WriteString(compositeKeyParts[11])
		buffer1.WriteString("\"")
		buffer1.WriteString("}")

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(compositeKeyParts[2])
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString("\"")
		buffer.WriteString(buffer1.String())
		buffer.WriteString("\"")
		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("-  queryCertificateBasedOnName returning:\n   %s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func createIndexHelper(stub shim.ChaincodeStubInterface, certificate *Certificate) error {
	var err error = nil

	for queryKey, indexName := range CerfificationQueryMap {
		if queryKey == "PartnerName" {
			err = createIndex(stub, indexName, []string{certificate.PartnerName, certificate.CertificateName, certificate.CertificateHash, certificate.Contacts, certificate.Mobile, certificate.Email, certificate.CertificateType, certificate.PassingDate, certificate.ExpiryDate, certificate.CertificateStatus, certificate.Participant, certificate.Score})
		} else if queryKey == "CertificateName" {
			err = createIndex(stub, indexName, []string{certificate.CertificateName, certificate.PartnerName, certificate.CertificateHash, certificate.Contacts, certificate.Mobile, certificate.Email, certificate.CertificateType, certificate.PassingDate, certificate.ExpiryDate, certificate.CertificateStatus, certificate.Participant, certificate.Score})
		}
	}

	return err
}

// ===============================================
// createIndex - create search index for ledger
// ===============================================
func createIndex(stub shim.ChaincodeStubInterface, indexName string, attributes []string) error {
	fmt.Println("- start create index")
	var err error
	indexKey, err := stub.CreateCompositeKey(indexName, attributes)
	if err != nil {
		return err
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of object.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(indexKey, value)

	fmt.Println("- end create index")
	return nil
}

func deleteIndexHelper(stub shim.ChaincodeStubInterface, certificate *Certificate) error {
	var err error = nil

	for queryKey, indexName := range CerfificationQueryMap {
		if queryKey == "PartnerName" {
			err = deleteIndex(stub, indexName, []string{certificate.PartnerName, certificate.CertificateName, certificate.CertificateHash, certificate.Contacts, certificate.Mobile, certificate.Email, certificate.CertificateType, certificate.PassingDate, certificate.ExpiryDate, certificate.CertificateStatus, certificate.Participant, certificate.Score})
		} else if queryKey == "CertificateName" {
			err = deleteIndex(stub, indexName, []string{certificate.CertificateName, certificate.PartnerName, certificate.CertificateHash, certificate.Contacts, certificate.Mobile, certificate.Email, certificate.CertificateType, certificate.PassingDate, certificate.ExpiryDate, certificate.CertificateStatus, certificate.Participant, certificate.Score})
		}
	}

	return err
}

func deleteIndex(stub shim.ChaincodeStubInterface, indexName string, attributes []string) error {
	fmt.Println("- start delete index")
	var err error
	indexKey, err := stub.CreateCompositeKey(indexName, attributes)
	if err != nil {
		return err
	}
	//  Delete index by key
	stub.DelState(indexKey)

	fmt.Println("- end delete index")
	return nil
}
