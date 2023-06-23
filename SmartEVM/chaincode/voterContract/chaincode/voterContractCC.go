package chaincode

import(
	"encoding/json"
	"fmt"
	"strconv"
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const Abstained="_Abstained_"

type SmartContract struct {
    contractapi.Contract
    IsAnonymous bool
    Initialized bool
    Abstainable bool
    SingleChoce bool
}


// NOTE: Write the asset properties in CAMEL CASE- otherwise, chaincode will not get deployed 
type Ballot struct {
	CastedBy  string    `json:"CastedBy"`
	Timestamp time.Time `json:"Timestamp"`
}
type VotableOption struct {
	VotableId string   `json:"VotableId"`
	Ballots   []Ballot `json:"Ballots"`
}

//Init Ledger
func (s *SmartContract) IntiLedger(ctx contractapi.TransactionContextInterface) {
	//TODO: values must come from chaincode deployment args
	s.Initialized = true
	s.IsAnonymous = true
	s.Abstainable = true
	s.SingleChoice = true
	if a.Abstainable {
		err := s.addVotableOption(ctx, Abstained)
		if err != nil {
			s.Initialized = false
			log.Fatal(err)
		}
	}
}

//Is Votable Option Exists?
func (s *SmartContract) isVotableOptionExists(ctx contractapi.TransactionContextInterface, votableId string) (bool, error) {
	if s.Initialized != true {
		return false,  fmt.Errorf("Ledger not inialized!")
	}
	optionBytes, err := ctx.GetStub().GetState(votableId)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", votableId, err)
	}

	return optionBytes != nil, nil
}

// function to add Votable Option 
func (s *SmartContract) addVotableOption(ctx contractapi.TransactionContextInterface, votableId string) error{

	if s.Initialized != true {
		return false,  fmt.Errorf("Ledger not inialized!")
	}
    	fmt.Printf("Adding new votable option: %s\n", votableId)
	//checking if option already added
	optionExists,err:= s.isVotableOptionExists(ctx, votableId)
	if err!=nil{
		return fmt.Errorf("error in checking whether asset exists or not: %v", err)
	}
	if optionExists {
		return fmt.Errorf("This votable option is already exist!")
	}

	// if the option does not exist
	votableOption := VotableOption{
		VotableId: votableId,
		Ballots: []Ballot
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
func (s *SmartContract) ReadOption(ctx contractapi.TransactionContextInterface, votableId string) (*VotableOption, error){
	if s.Initialized != true {
		return false,  fmt.Errorf("Ledger not inialized!")
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
func (s *SmartContract) castVote(ctx contractapi.TransactionContextInterface,userId string, votableIds []string) error {
	if s.Initialized != true {
		return false,  fmt.Errorf("Ledger not inialized!")
	}
	
	if len(votableIds) == 0 {
		if s.IsAbstainable {
			return castVote(str, {Abstained})
		}
		return fmt.Error("Number of votes to be casted can not be zero.")
	}
	if ( len(votableIds) > 1 && s.SingleChoice ) {
		return fmt.Error("Number of votes can not be more than one, in case of single choice voting.")
	}
	//Get all current voting options
	fmt.Println("Fetching all possible votable options ...")
	options, err := s.GetAllOptions(ctx)
	if err != nil {
		fmt.Println("Could not get all voting options : %v", err)
		return err
	}
	
	var updatedOptions []VotableOptions

	for i, votableId := range options {
		fmt.Println(i, "->", votableId) 
		votableOption, err := s.ReadOption(ctx, votableId) // asset is read
		if err != nil {
			fmt.Println("Could not retrieve option %s : %v", votableId, err)
			continue
		}

		for j, ballot := votableOption.Ballots {
			fmt.Println("   Ballot ", j, "->", ballot.Timestamp, ";", ballot.CastedBy)
			if ballot.CastedBy == userId {
				return fmt.Error("User %s already casted the vote!", userId)
			}
		}
		if contains(votableIds, votablId) {
			append(votableId.Ballots, Ballot {CastedBy: s.IsAnonymous?"Anonymous":userId, Timestamp: time.Now})
			append(updatedOptions, votableId)
		}
	}
	for i, updated := range updatedOptions {
		fmt.Println("Registering vote for %s : %s", userId, updated.VotableId)
		updatedJSON, err := json.Marshal(updated)
		if err != nil {
			return err
		}
		updatestate_err := ctx.GetStub().PutState(updated.VotableId, updatedJSON)
	}
	return nil
   }
}
// GetAllAssets returns all voting options found in world state
func (s *SmartContract) GetAllOptions(ctx contractapi.TransactionContextInterface) ([]*VotingOption, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*VotingOption
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset VotingOption
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
