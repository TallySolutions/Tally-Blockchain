package chaincode

import(
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

type Company struct {
	companyName string `json:companyName`
	licenseId string `json:licenseId`
	// address
	// reg number
	// gst number
}

type TallyScoreAsset struct {
    companyLicenseId string `json:companyLicenseId`
	score uint `json:score`
}

//function to check whether the company has already been registered
func (s *SmartContract) companyAssetExists(ctx contractapi.TransactionContextInterface, licenseId string) (bool, error) {

    companyAssetJSON, err := ctx.GetStub().GetState(licenseId)
    if err != nil {
    	return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return companyAssetJSON != nil, nil
}

// function to register company and initialize it's score to 500 
func (s *SmartContract) RegisterCompany(ctx contractapi.TransactionContextInterface, licenseId string) error{

	//checking if licenseID is valid
	companyAssetExists,err:= s.companyAssetExists(ctx, licenseId)
	if err!=nil{
		return fmt.Errorf("error in checking whether asset exists: %v", err)
	}
	if companyAssetExists {
		return fmt.Errorf("This company already exists!")
	}

	// if the company is unregistered
	companyScoreAsset := TallyScoreAsset{
		companyLicenseId: licenseId,
		score: 500,
	}
	companyScoreAssetJSON, err := json.Marshal(companyScoreAsset)
    if err != nil {
        return err
    }

	putStateErr := ctx.GetStub().PutState(licenseId, companyScoreAssetJSON) // new state added to the tallyscore ledger
    fmt.Printf("Asset creation returned : %s\n", putStateErr)
    return putStateErr

} 

// function to unregister a company (deleting it's score asset)
func (s *SmartContract) UnregisterCompany(ctx contractapi.TransactionContextInterface, licenseId string) error{
	//checking if licenseID is valid
	var sumOfDigits int
	for _, charDigit:= range licenseId{
		digit:= int(charDigit- '0')
		sumOfDigits+= digit
	}
	if sumOfDigits%9 !=0{
		return fmt.Errorf("Invalid license ID")
	}

	exists, err := s.companyAssetExists(ctx, licenseId)
        if err != nil {
            return err
        }
        if !exists {
            return fmt.Errorf("the asset %s does not exist", licenseId)
        }
	delStateOp:= ctx.GetStub().DelState(licenseId)
    fmt.Printf("Message received on deletion: %s", delStateOp)
    return nil
}

// function to read companyasset
func (s *SmartContract) ReadCompanyAsset(ctx contractapi.TransactionContextInterface, licenseID string) (*TallyScoreAsset, error){
	companyScoreAssetJSON, err := ctx.GetStub().GetState(licenseID)
    if err != nil {
    	return nil, fmt.Errorf("Failed to read from world state: %v", err)
    }
    if companyScoreAssetJSON == nil {
    	return nil, fmt.Errorf("The company with ID %s is not registered.", licenseID)
    }

    var companyScoreAsset TallyScoreAsset
    err = json.Unmarshal(companyScoreAssetJSON, &companyScoreAsset)
    if err != nil {
    	return nil, err
	}

	return &companyScoreAsset, nil
}

// function to increase tallyScore of a company
func (s *SmartContract) IncreaseScore(ctx contractapi.TransactionContextInterface, licenseID string, incrementValue string) (*TallyScoreAsset, error) {
    companyAssetRead, err := s.ReadCompanyAsset(ctx, licenseID) // asset is read
    if err != nil {
    	return nil, err
    }

    intermediateUpdateval, err := strconv.ParseUint(incrementValue, 10, 32)
    if err !=nil {
    	fmt.Println(err)
    }
	newScore:= uint(companyAssetRead.score) + ((1000- companyAssetRead.score) * uint(intermediateUpdateval))/100
    if newScore > 1000 {
    	return nil, fmt.Errorf("You cannot have a value more than 1000.")
    }

    // overwriting original asset with new value
    companyAsset := TallyScoreAsset {
        companyLicenseId: licenseID,
		score: newScore,
    }
    companyAssetJSON, err := json.Marshal(companyAsset)
    if err != nil {
    	return nil, err
	}

	updatestate_err := ctx.GetStub().PutState(licenseID, companyAssetJSON)
	fmt.Printf("Increasing company asset score returned the following: %s ", updatestate_err)
	return &companyAsset, nil
}

// function to decrease tallyScore of a company
func (s *SmartContract) DecreaseScore(ctx contractapi.TransactionContextInterface, licenseID string, decrementValue string) (*TallyScoreAsset, error) {
    companyAssetRead, err := s.ReadCompanyAsset(ctx, licenseID) // asset is read
    if err != nil {
    	return nil, err
    }

    intermediateUpdateval, err := strconv.ParseUint(decrementValue, 10, 32)
    if err !=nil {
    	fmt.Println(err)
    }
	newScore:= uint(companyAssetRead.score) - ((1000- companyAssetRead.score) * uint(intermediateUpdateval))/100
    if newScore < 0 {
    	return nil, fmt.Errorf("You cannot have a value lesser than 0.")
    }

    // overwriting original asset with new value
    companyAsset := TallyScoreAsset {
        companyLicenseId: licenseID,
		score: newScore,
    }
    companyAssetJSON, err := json.Marshal(companyAsset)
    if err != nil {
    	return nil, err
	}

	updatestate_err := ctx.GetStub().PutState(licenseID, companyAssetJSON)
	fmt.Printf("Decreasing company asset score returned the following: %s ", updatestate_err)
	return &companyAsset, nil
}