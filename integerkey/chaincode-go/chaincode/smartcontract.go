package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	Name  string `json:"Name"`
	Value uint   `json:"Value"`
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, Name string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(Name)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, Name string) error {
	exists, err := s.AssetExists(ctx, Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", Name)
	}

	asset := Asset{
		Name:  Name,
		Value: 0,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err := ctx.GetStub().PutState(Name, assetJSON)

	fmt.Println("Asset creation returnde : %s", err)

	return err
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
func (s *SmartContract) IncreaseAsset(ctx contractapi.TransactionContextInterface, Name string, incrementValue uint) error {
	// exists, err := s.AssetExists(ctx, Name)

	asset_read, err := s.ReadAsset(ctx, Name)
	if err != nil {
		return err
	}
	// use GetState()
	newValue := uint(asset_read.Value) + incrementValue

	if newValue > 20 {
		return fmt.Errorf("You cannot have a value more than 20.")
	}

	// overwriting original asset with new value
	asset := Asset{
		Name:  Name,
		Value: newValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(Name, assetJSON)
}

func (s *SmartContract) DecreaseAsset(ctx contractapi.TransactionContextInterface, Name string, decrementValue uint) error {
	asset_read, err := s.ReadAsset(ctx, Name)
	if err != nil {
		return err
	}
	newValue := uint(asset_read.Value) - decrementValue

	// overwriting original asset with new value
	asset := Asset{
		Name:  Name,
		Value: newValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(Name, assetJSON)
}


