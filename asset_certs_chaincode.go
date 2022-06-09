package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Chaincode struct {
	contractapi.Contract
}

//aCert stands for academic certificates
type AcademicCert struct {
	DocType        string   `jason:"docType"`
	ACertID        string   `jason:"aCertID"`
	StudentID      string   `jason:"studentID"`
	StudentName    string   `jason:"studentName"`
	Degree         string   `jason:"degree"`
	GraduationDate string   `jason:"graduationDate"`
	Transcript     []string `jason:"transcript"`
}

//cCert stands for extra-curricular certificates
type ExtraCurricularCert struct {
	DocType      string   `jason:"docType"`
	CCertID      string   `jason:"cCertID"`
	StudentID    string   `jason:"studentID"`
	StudentName  string   `jason:"studentName"`
	Achievements []string `jason:"achievements"`
}

// issueAcaCert initializes a new academic certificate on the blockchain
func (t *Chaincode) IssueAcaCert(ctx contractapi.TransactionContextInterface, aCertID, studentID string, studentName string, degree string, gradDate string, transcript []string) error {
	exists, err := t.AssetExists(ctx, aCertID)
	if err != nil {
		return fmt.Errorf("failed to get academic cert: %v", err)
	}
	if exists {
		return fmt.Errorf("academic cert already exists: %s", aCertID)
	}

	aCert := &AcademicCert{
		DocType:        "aCert",
		ACertID:        aCertID,
		StudentID:      studentID,
		StudentName:    studentName,
		Degree:         degree,
		GraduationDate: gradDate,
		Transcript:     transcript,
	}
	aCertBytes, err := json.Marshal(aCert)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(aCertID, aCertBytes)
	return err
}

//issueCert initializes a new extracurricular certificate on the blockchain
func (t *Chaincode) IssueCurrCert(ctx contractapi.TransactionContextInterface, cCertID, studentID string, studentName string, achievements []string) error {
	exists, err := t.AssetExists(ctx, cCertID)
	if err != nil {
		return fmt.Errorf("failed to get extracurricular cert: %v", err)

	}
	if exists {
		return fmt.Errorf("extracurricular cert already exists: %s", cCertID)
	}

	cCert := &ExtraCurricularCert{
		DocType:      "cCert",
		CCertID:      cCertID,
		StudentID:    studentID,
		StudentName:  studentName,
		Achievements: achievements,
	}
	cCertBytes, err := json.Marshal(cCert)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(cCertID, cCertBytes)
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

//ReadCurrCert retrieves a extracurricular certificate from the ledger
func (t *Chaincode) ReadCurrCert(ctx contractapi.TransactionContextInterface, cCertID string) (*ExtraCurricularCert, error) {
	cCertBytes, err := ctx.GetStub().GetState(cCertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get extracurricular cert %s: %v", cCertID, err)
	}
	if cCertBytes == nil {
		return nil, fmt.Errorf("extracurricular cert %s does not exists", cCertID)
	}

	var cCert ExtraCurricularCert
	err = json.Unmarshal(cCertBytes, &cCert)
	if err != nil {
		return nil, err
	}

	return &cCert, nil
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator that returns an array of academic certs
func constructQueryResponseFromIterator_AcaCert(resultsIterator shim.StateQueryIteratorInterface) ([]*AcademicCert, error) {
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

// QueryAcaCertByStudentID queries for academic certificates based on the studentID.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (studentID).
func (t *Chaincode) QueryAcaCertByStudentID(ctx contractapi.TransactionContextInterface, studentID string) ([]*AcademicCert, error) {
	queryString := fmt.Sprintf(`{"selector":{"DocType":"aCert","StudentID":"%s"}}`, studentID)
	return getQueryResultForQueryString_AcaCert(ctx, queryString)
}

// getQueryResultForQueryString_AcaCert executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString_AcaCert(ctx contractapi.TransactionContextInterface, queryString string) ([]*AcademicCert, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator_AcaCert(resultsIterator)
}

// constructQueryResponseFromIterator_CurrCert constructs a slice of assets from the resultsIterator that returns an array of extracurricular certifications
func constructQueryResponseFromIterator_CurrCert(resultsIterator shim.StateQueryIteratorInterface) ([]*ExtraCurricularCert, error) {
	var cCerts []*ExtraCurricularCert
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var cCert ExtraCurricularCert
		err = json.Unmarshal(queryResult.Value, &cCert)
		if err != nil {
			return nil, err
		}
		cCerts = append(cCerts, &cCert)
	}
	return cCerts, nil
}

// QueryCurrCertByStudentID queries for extracurricular certificates based on the studentID.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (studentID).
func (t *Chaincode) QueryCurrCertByStudentID(ctx contractapi.TransactionContextInterface, studentID string) ([]*ExtraCurricularCert, error) {
	queryString := fmt.Sprintf(`{"selector":{"DocType":"cCert","StudentID":"%s"}}`, studentID)
	return getQueryResultForQueryString_CurrCert(ctx, queryString)
}

// getQueryResultForQueryString_CurrCert executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString_CurrCert(ctx contractapi.TransactionContextInterface, queryString string) ([]*ExtraCurricularCert, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator_CurrCert(resultsIterator)
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
		{DocType: "aCert", ACertID: "aCert1", StudentID: "SWE1904873", StudentName: "Loo Yong Jun", Degree: "Bachelor of Chemical Engineering with Honours", GraduationDate: "23 April 2020", Transcript: []string{"Introduction to software engineering, GPA : 4.0", "Computing architecture, GPA : 3.3"}},
		{DocType: "aCert", ACertID: "aCert2", StudentID: "SWE1909886", StudentName: "Shane Ho Ken Wae", Degree: "Bachelor of Science in Mathematics and Applied Mathematics with Honours", GraduationDate: "20 September 2020", Transcript: []string{"Introduction to software engineering, GPA : 2.0", "International Business, GPA : 2.0"}},
	}

	cCerts := []ExtraCurricularCert{
		{DocType: "cCert", CCertID: "cCert1", StudentID: "SWE1904873", StudentName: "Loo Yong Jun", Achievements: []string{"Club: Martial Arts Club", "Position: Vice President", "Service: April 2021 - April 2022"}},
		{DocType: "cCert", CCertID: "cCert2", StudentID: "SWE1904873", StudentName: "Loo Yong Jun", Achievements: []string{"Event: X-Tech Hackathon", "Award: Participation", "Date: 22 September 2020 - 24 September 2022", "Organizer: Tech Club (X-Tech)"}},
		{DocType: "cCert", CCertID: "cCert3", StudentID: "SWE1909886", StudentName: "Shane Ho Ken Wae", Achievements: []string{"Club: Music Club", "Position: General Affairs", "Service: April 2019 - April 2020"}},
	}

	for _, acert := range aCerts {
		err := t.IssueAcaCert(ctx, acert.ACertID, acert.StudentID, acert.StudentName, acert.Degree, acert.GraduationDate, acert.Transcript)
		if err != nil {
			return err
		}
	}

	for _, ccert := range cCerts {
		err := t.IssueCurrCert(ctx, ccert.CCertID, ccert.StudentID, ccert.StudentName, ccert.Achievements)
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
