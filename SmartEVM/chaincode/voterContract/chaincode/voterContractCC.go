package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const Abstained = "_Abstained_"

type SmartContract struct {
	contractapi.Contract
	IsAnonymous  bool
	Initialized  bool
	Abstainable  bool
	SingleChoice bool
}

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type DateTime time.Time

func (t DateTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("RFC3339"))
	return []byte(stamp), nil
}

// NOTE: Write the asset properties in CAMEL CASE- otherwise, chaincode will not get deployed
type Ballot struct {
	CastedBy  string   `json:"CastedBy"`
	Timestamp DateTime `json:"Timestamp"`
}
type VotableOption struct {
	VotableId string   `json:"VotableId"`
	Ballots   []Ballot `json:"Ballots"`
}

// Utility function
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Init Ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	//TODO: values must come from chaincode deployment args
	s.Initialized = true
	s.IsAnonymous = true
	s.Abstainable = true
	s.SingleChoice = true
	if s.Abstainable {
		err := s.AddVotableOption(ctx, Abstained)
		if err != nil {
			s.Initialized = false
			return err
		}
	}
	return nil
}

// Is Votable Option Exists?
func (s *SmartContract) isVotableOptionExists(ctx contractapi.TransactionContextInterface, votableId string) (bool, error) {
	if s.Initialized != true {
		return false, fmt.Errorf("Ledger not inialized!")
	}
	optionBytes, err := ctx.GetStub().GetState(votableId)
	if err != nil {
		return false, err
	}

	return optionBytes != nil, nil
}

// function to add Votable Option
func (s *SmartContract) AddVotableOption(ctx contractapi.TransactionContextInterface, votableId string) error {

	if s.Initialized != true {
		return fmt.Errorf("Ledger not inialized!")
	}
	fmt.Printf("Adding new votable option: %s\n", votableId)
	//checking if option already added
	optionExists, err := s.isVotableOptionExists(ctx, votableId)
	if err != nil {
		return err
	}
	if optionExists {
		return fmt.Errorf("This votable option is already exist!")
	}

	// if the option does not exist
	var ballots []Ballot
	votableOption := VotableOption{
		VotableId: votableId,
		Ballots:   ballots,
	}
	votableOptionJSON, err := json.Marshal(votableOption)
	if err != nil {
		return err
	}

	fmt.Printf("Creating new asset for this votable id in voting ledger: %s\n", votableId)
	putStateErr := ctx.GetStub().PutState(votableId, votableOptionJSON) // new state added to the voting ledger
	return putStateErr

}

// function to read vote option
func (s *SmartContract) ReadOption(ctx contractapi.TransactionContextInterface, votableId string) (*VotableOption, error) {
	if s.Initialized != true {
		return nil, fmt.Errorf("Ledger not inialized!")
	}
	votableOptionJSON, err := ctx.GetStub().GetState(votableId)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if votableOptionJSON == nil {
		return nil, fmt.Errorf("The votable id %s is not present.", votableId)
	}

	var votableOption VotableOption
	err = json.Unmarshal(votableOptionJSON, &votableOption)
	if err != nil {
		return nil, err
	}

	return &votableOption, nil
}

// function to cast vote
func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, userId string, votableIds []string) error {
	if s.Initialized != true {
		return fmt.Errorf("Ledger not inialized!")
	}

	if len(votableIds) == 0 {
		if s.Abstainable {
			return s.CastVote(ctx, userId, []string{Abstained})
		}
		return fmt.Errorf("Number of votes to be casted can not be zero.")
	}
	if len(votableIds) > 1 && s.SingleChoice {
		return fmt.Errorf("Number of votes can not be more than one, in case of single choice voting.")
	}
	//Get all current voting options
	fmt.Println("Fetching all possible votable options ...")
	options, err := s.GetAllOptions(ctx)
	if err != nil {
		fmt.Println("Could not get all voting options : ", err)
		return err
	}

	var updatedOptions []VotableOption

	for i, votableOption := range options {
		fmt.Println(i, "->", votableOption.VotableId)

		for j, ballot := range votableOption.Ballots {
			fmt.Println("   Ballot ", j, "->", ballot.Timestamp, ";", ballot.CastedBy)
			if ballot.CastedBy == userId {
				return fmt.Errorf("User %s already casted the vote!", userId)
			}
		}
		if contains(votableIds, votableOption.VotableId) {
			castedBy := userId
			if s.IsAnonymous {
				castedBy = "Anonymous"
			}
			ballot := Ballot{CastedBy: castedBy, Timestamp: DateTime(time.Now())}
			votableOption.Ballots = append(votableOption.Ballots, ballot)
			updatedOptions = append(updatedOptions, *votableOption)
		}
	}
	for i, updated := range updatedOptions {
		fmt.Printf("(%d) Registering vote for %s : %s\n", i, userId, updated.VotableId)
		updatedJSON, err := json.Marshal(updated)
		if err != nil {
			return err
		}
		updatestate_err := ctx.GetStub().PutState(updated.VotableId, updatedJSON)
		if updatestate_err != nil {
			return updatestate_err
		}
	}
	return nil
}

// GetAllAssets returns all voting options found in world state
func (s *SmartContract) GetAllOptions(ctx contractapi.TransactionContextInterface) ([]*VotableOption, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}

	if resultsIterator == nil {
		return nil, fmt.Errorf("No result iterator found!")
	}

	defer resultsIterator.Close()

	var assets []*VotableOption
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		fmt.Println(queryResponse.Value)
		var asset VotableOption
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
