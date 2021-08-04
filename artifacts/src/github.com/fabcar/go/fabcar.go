package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
	contractapi.Contract
}

type Patient struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ID        string `json:"id"`
	//Records   []Record `json:"records"`
}

// type Doctor struct {
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// 	ID        string `json:"id"`
// 	Hospital  string `json:"hospital"`
// }

// type Record struct {
// 	Owner       Patient `json:"owner"`
// 	ID          string  `json:"id"`
// 	Information string  `json:"information"`
// }

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Patient
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Init()", fcn, params)
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	var result []byte
	var err error
	if fcn == "CreatePatient" {
		result, err = s.CreatePatient(stub, params)
	} else if fcn == "GetPatient" {
		result, err = s.GetPatient(stub, params)
	} else if fcn == "GetAllPatients" {
		result, err = s.GetAllPatients(stub)
	}
	// } else if fcn == "AddRecordToPatient" {
	// 	result, err = s.AddRecordToPatient(stub, params)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func (s *SmartContract) CreatePatient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Create Patient executed")

	if len(args) != 3 {
		return nil, fmt.Errorf("failed to create Patient: The number of arguments is incorrect")

	}

	//Create new Patient
	fmt.Printf("args: %s", args[2])

	var newPatient = Patient{FirstName: args[0], LastName: args[1], ID: args[2]}

	newPatientAsBytes, err := json.Marshal(newPatient)
	if err != nil {
		return nil, fmt.Errorf("failed to create Patient")
	}

	stub.PutState(args[2], newPatientAsBytes)
	return newPatientAsBytes, nil
}

func (s *SmartContract) GetPatient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("failed to get Patient: The number of arguments is incorrect")

	}

	patientAsBytes, err := stub.GetState(args[2])

	if err != nil {
		return nil, fmt.Errorf("failed to get Patient %s", args[2])
	}

	if patientAsBytes == nil {
		return nil, fmt.Errorf("failed to get Patient %s: It doet not exists", args[2])
	}

	return patientAsBytes, nil
}

func (s *SmartContract) GetAllPatients(stub shim.ChaincodeStubInterface) ([]byte, error) {
	iterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get all Patients")
	}

	defer iterator.Close()

	var buffer bytes.Buffer
	first := true
	buffer.WriteString("[")

	for iterator.HasNext() {
		next, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get all Patients")
		}

		if !first {
			buffer.WriteString(", ")
		} else {
			first = false
		}

		buffer.WriteString("{\"Key\": \"")
		buffer.WriteString(next.Key)
		buffer.WriteString("\", \"Value\": ")
		buffer.Write(next.Value)
		buffer.WriteString("}")
	}

	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
