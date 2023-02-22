package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	Name  string `json:"Name"`
	Value uint   `json:"Value"`
}

// function that takes input as context of transaction and the name of the key, returns boolean value that implies whether the asset exists or not, otherwise- an error
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, Name string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(Name)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// function to create an asset. Input= transaction context, name of the key to be created. Creates new asset if an asset with the name given does not exist
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, Name string) error {

	exists, err := s.AssetExists(ctx, Name) // exists-> boolean value, err-> can be nil or the error, if present

	fmt.Printf("Asset exists returned : %t, %s\n", exists, err)

	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", Name)
	}

	asset := Asset{ //creation of asset
		Name:  Name,
		Value: 0,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	state_err := ctx.GetStub().PutState(Name, assetJSON) // new state added

	fmt.Printf("Asset creation returned : %s\n", state_err)

	return state_err
}

// ReadAsset returns the asset stored in the world state with given Name.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, Name string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(Name)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", Name)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([] *Asset, error){

	iteratorVar, err := ctx.GetStub().GetStateByRange("","")   // TRY RANGE PARAMETERS , other getstateby.... (rows etc.)
	if err !=nil{
		return nil, err
	}
	defer iteratorVar.Close()


	var assets []*Asset

	for iteratorVar.HasNext() {
		queryResponse, err := iteratorVar.Next()
		if err != nil {
		  return nil, err
		}
	
		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
		  return nil, err
		}
		assets = append(assets, &asset)
	  }
	
	  return assets, nil

}



// IncreaseAsset increases the value of the asset by the specified value- with certain limits
func (s *SmartContract) IncreaseAsset(ctx contractapi.TransactionContextInterface, Name string, incrementValue string) (*Asset, error) {
	// NOTE: incrementValue is a string because SubmitTransaction accepts string parameters as input parameters
	asset_read, err := s.ReadAsset(ctx, Name) // asset is read
	if err != nil {
		return nil, err
	}

	intermediateUpdateval, err := strconv.ParseUint(incrementValue, 10, 32)  
	if err !=nil{
			fmt.Println(err)
	}
	incrementValueuInt := uint(intermediateUpdateval)
	newValue := uint(asset_read.Value) + incrementValueuInt

	if newValue > 20 {
		return nil, fmt.Errorf("You cannot have a value more than 20.")
	}

	// overwriting original asset with new value
	asset := Asset{
		Name:  Name,
		Value: newValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	updatestate_err := ctx.GetStub().PutState(Name, assetJSON)
	fmt.Printf("Increasing asset value returned the following: %s ", updatestate_err)

	return &asset, nil
}

// DecreaseAsset decreases the value of the asset by the specified value
func (s *SmartContract) DecreaseAsset(ctx contractapi.TransactionContextInterface, Name string, decrementValue string) (*Asset, error) {
	asset_read, err := s.ReadAsset(ctx, Name)
	if err != nil {
		return nil, err
	}

	intermediateval, err := strconv.ParseUint(decrementValue, 10, 32)
	if err !=nil{
			fmt.Println(err)
	}
	decrementValueuInt := uint(intermediateval)
	newValue := uint(asset_read.Value) - decrementValueuInt

	// overwriting original asset with new value
	asset := Asset{
		Name:  Name,
		Value: newValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	updatestate_Err := ctx.GetStub().PutState(Name, assetJSON)
	fmt.Printf("After decreasing asset value: %s", updatestate_Err)

	return &asset , nil
}