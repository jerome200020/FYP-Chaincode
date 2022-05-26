package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "uniCert~studentID"

type Chaincode struct {
	contractapi.Contract
}

//aCert stands for academic certificates
type AcademicCert struct {
	DocType     string   `jason:"docType"`
	aCertID     string   `jason:"aCertID"`
	studentID   string   `jason:"studentID"`
	studentName string   `jason:"studentName"`
	transcript  []string `jason:"transcript"`
}

// issueCert initializes a new academic certificate on the blockchain
func (t *Chaincode) IssueAcaCert(ctx contractapi.TransactionContextInterface, aCertID, studentID string, studentName string, transcript []string) error {
	exists, err := t.AssetExists(ctx, aCertID)
	if err != nil {
		return fmt.Errorf("failed to get academic cert: %v", err)
	}
	if exists {
		return fmt.Errorf("academic cert already exists: %s", aCertID)
	}

	aCert := &AcademicCert{
		DocType:     "aCert",
		aCertID:     aCertID,
		studentID:   studentID,
		studentName: studentName,
		transcript:  transcript,
	}
	aCertBytes, err := json.Marshal(aCert)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(aCertID, aCertBytes)
	return err
}

// ReadAcaCert retrieves a academic certificate from the ledger
func (t *Chaincode) ReadAcaCert(ctx contractapi.TransactionContextInterface, aCertID string) (*AcademicCert, error) {
	aCertBytes, err := ctx.GetStub().GetState(aCertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get academic cert %s: %v", aCertID, err)
	}
	if aCertBytes == nil {
		return nil, fmt.Errorf("academic cert %s does not exists", aCertID)
	}

	var aCert AcademicCert
	err = json.Unmarshal(aCertBytes, &aCert)
	if err != nil {
		return nil, err
	}

	return &aCert, nil
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*AcademicCert, error) {
	var aCerts []*AcademicCert
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var aCert AcademicCert
		err = json.Unmarshal(queryResult.Value, &aCert)
		if err != nil {
			return nil, err
		}
		aCerts = append(aCerts, &aCert)
	}

	return aCerts, nil
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*AcademicCert, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// QueryAssetsByOwner queries for academic certificates based on the studentID.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (studentID).
func (t *Chaincode) QueryAcaCertByStudentID(ctx contractapi.TransactionContextInterface, studentID string) ([]*AcademicCert, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"aCert","studentID":"%s"}}`, studentID)
	return getQueryResultForQueryString(ctx, queryString)
}

// AssetExists returns true when asset with given ID exists in the ledger.
func (t *Chaincode) AssetExists(ctx contractapi.TransactionContextInterface, certID string) (bool, error) {
	certBytes, err := ctx.GetStub().GetState(certID)
	if err != nil {
		return false, fmt.Errorf("failed to read cert %s from world state. %v", certID, err)
	}

	return certBytes != nil, nil
}

func (t *Chaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	aCerts := []AcademicCert{
		{DocType: "aCert", aCertID: "aCert1", studentID: "SWE1904873", studentName: "Loo Yong Jun",
			transcript: []string{"Introduction to software engineering, GPA : 4.0",
				"Computing architecture, GPA : 3.3"}},
		{DocType: "aCert", aCertID: "aCert2", studentID: "SWE1909886", studentName: "Loo Ken Wae",
			transcript: []string{"Introduction to software engineering, GPA : 2.0",
				"International Business, GPA : 2.0"}},
	}

	for _, cert := range aCerts {
		err := t.IssueAcaCert(ctx, cert.aCertID, cert.studentID, cert.studentName, cert.transcript)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&Chaincode{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
