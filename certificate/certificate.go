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
	} else if function == "queryCertificateBasedOnParterName" {
		return s.queryCertificateBasedOnParterName(stub, args)
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
	indexName := "partnername~all"
	err = createIndex(stub, indexName, []string{certificate.PartnerName, certificate.CertificateHash, certificate.Contacts, certificate.Mobile, certificate.Email, certificate.CertificateType, certificate.CertificateName, certificate.PassingDate, certificate.ExpiryDate, certificate.CertificateStatus, certificate.Participant, certificate.Score})
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
	indexName := "partnername~all"
	err = deleteIndex(stub, indexName, []string{certificate.PartnerName, certificate.CertificateHash, certificate.Contacts, certificate.Mobile, certificate.Email, certificate.CertificateType, certificate.CertificateName, certificate.PassingDate, certificate.ExpiryDate, certificate.CertificateStatus, certificate.Participant, certificate.Score})
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
	indexName := "partnername~all"
	err := deleteIndex(stub, indexName, []string{certificate.PartnerName, certificate.CertificateHash, certificate.Contacts, certificate.Mobile, certificate.Email, certificate.CertificateType, certificate.CertificateName, certificate.PassingDate, certificate.ExpiryDate, certificate.CertificateStatus, certificate.Participant, certificate.Score})
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
	err = createIndex(stub, indexName, []string{certificate.PartnerName, certificate.CertificateHash, certificate.Contacts, certificate.Mobile, certificate.Email, certificate.CertificateType, certificate.CertificateName, certificate.PassingDate, certificate.ExpiryDate, certificate.CertificateStatus, certificate.Participant, certificate.Score})
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

func (s *SmartContract) queryCertificateBasedOnParterName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	partnerName := args[0]
	fmt.Println("- start  queryCertificateBasedOnParterName", partnerName)

	certificateResultsIterator, err := stub.GetStateByPartialCompositeKey("partnername~all", []string{partnerName})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer certificateResultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	var i int
	for i = 0; certificateResultsIterator.HasNext(); i++ {
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
		returnedPartnerName := compositeKeyParts[0]
		returnedCertificateHash := compositeKeyParts[1]
		returnedContacts := compositeKeyParts[2]
		returnedMobile := compositeKeyParts[3]
		returnedEmail := compositeKeyParts[4]
		returnedCertificateType := compositeKeyParts[5]
		returnedCertificateName := compositeKeyParts[6]
		returnedPassingDate := compositeKeyParts[7]
		returnedExpiryDate := compositeKeyParts[8]
		returnedCertificateStatus := compositeKeyParts[9]
		returnedParticipant := compositeKeyParts[10]
		returnedScore := compositeKeyParts[11]
		fmt.Printf("- found a certificate record from index:%s partnerName:%s certificateHash:%s contacts:%s mobile:%s email:%s certificateType:%s certificateName:%s passingDate:%s expiryDate:%s certificateStatus:%s participant:%s score:%s\n",
			objectType, returnedPartnerName, returnedCertificateHash,
			returnedContacts, returnedMobile, returnedEmail, returnedCertificateType,
			returnedCertificateName, returnedPassingDate, returnedExpiryDate,
			returnedCertificateStatus, returnedParticipant, returnedScore)
		buffer.WriteString("{\"PartnerName\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedPartnerName)
		buffer.WriteString("\"")

		buffer.WriteString(", \"CertificateHash\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedCertificateHash)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Contacts\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedContacts)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Mobile\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedMobile)
		buffer.WriteString("\"")
		buffer.WriteString("}")

		buffer.WriteString(", \"Email\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedEmail)
		buffer.WriteString("\"")

		buffer.WriteString(", \"CertificateType\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedCertificateType)
		buffer.WriteString("\"")

		buffer.WriteString(", \"CertificateName\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedCertificateName)
		buffer.WriteString("\"")

		buffer.WriteString(", \"PassingDate\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedPassingDate)
		buffer.WriteString("\"")

		buffer.WriteString(", \"ExpiryDate\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedExpiryDate)
		buffer.WriteString("\"")

		buffer.WriteString(", \"CertificateStatus\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedCertificateStatus)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Participant\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedParticipant)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Score\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedScore)
		buffer.WriteString("\"")

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("-  queryCertificateBasedOnParterName returning:\n   %s\n", buffer.String())
	return shim.Success(buffer.Bytes())
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
