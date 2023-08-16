package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type ElectionConfig struct {
	IsAnonymous   bool `json:"anonymous"`
	IsSingle      bool `json:"single_choice"`
	IsAbstainable bool `json:"abstainable"`
}

type Option struct {
	ID    string `json:"id"`
	Votes int    `json:"votes"`
}

type Voter struct {
	ID string `json:"id"`
}

type Ballot struct {
	VoterID   string    `json:"voterID"`
	OptionIDs []string  `json:"optionIDs"`
	HasVoted  bool      `json:"hasVoted"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, isAnonymous bool, isSingle bool, isAbstainable bool) error {
	//create the election configuration
	electionConfig := ElectionConfig{
		IsAnonymous:   isAnonymous,
		IsSingle:      isSingle,
		IsAbstainable: isAbstainable,
	}
	//marshal and save the election configuration in the ledger
	electionConfigBytes, err := json.Marshal(electionConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal election config: %v", err)
	}

	err = ctx.GetStub().PutState("electionConfig", electionConfigBytes)
	if err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}

	// initial options to be registered in the ledger
	options := []Option{
		{ID: "1", Votes: 0},
		{ID: "2", Votes: 0},
		{ID: "3", Votes: 0},
	}

	// Loop through the options and register each one in the ledger
	for _, option := range options {
		optionBytes, err := json.Marshal(option)
		if err != nil {
			return fmt.Errorf("failed to marshal option: %v", err)
		}

		err = ctx.GetStub().PutState(option.ID, optionBytes)
		if err != nil {
			return fmt.Errorf("failed to put state: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) AddVoters(ctx contractapi.TransactionContextInterface, voterIDs []string) error {
	for _, voterID := range voterIDs {
		existingVoterJSON, err := ctx.GetStub().GetState(voterID)
		if err != nil {
			return fmt.Errorf("failed to read voter state from the ledger:%w", err)
		}
		if existingVoterJSON != nil {
			return fmt.Errorf("voter with ID %s already exists", voterID)
		}

		voter := Voter{
			ID: voterID,
		}

		voterJSON, err := json.Marshal(voter)
		if err != nil {
			return fmt.Errorf("failed to marshal voter JSON :%w", err)
		}

		err = ctx.GetStub().PutState(voterID, voterJSON)
		if err != nil {
			return fmt.Errorf("failed to insert voter: %v", err)
		}

		//Create an empty ballot for the voter with hasVoted set to false
		ballot := Ballot{
			VoterID:   voterID,
			OptionIDs: []string{},
			HasVoted:  false,
			Timestamp: time.Now(),
		}
		ballotJSON, err := json.Marshal(ballot)
		if err != nil {
			return fmt.Errorf("failed to marshal ballot JSON: %v", err)
		}
		err = ctx.GetStub().PutState(voterID+"ballot", ballotJSON)
		if err != nil {
			return fmt.Errorf("failed to insert ballot: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) RegisterOptions(ctx contractapi.TransactionContextInterface, optionIDs []string) error {
	for _, optionID := range optionIDs {
		// Check if the option with the given ID already exists in the ledger
		optionJSON, err := ctx.GetStub().GetState(optionID)
		if err != nil {
			return fmt.Errorf("failed to read from world state: %w", err)
		}
		if optionJSON != nil {
			return fmt.Errorf("option with ID %s already exists", optionID)
		}
		// Create a new option with initial votes set to 0
		option := Option{
			ID:    optionID,
			Votes: 0,
		}
		// Convert the option to JSON and save it to the ledger
		optionJSON, err = json.Marshal(option)
		if err != nil {
			return fmt.Errorf("failed to marshal option: %w", err)
		}
		err = ctx.GetStub().PutState(optionID, optionJSON)
		if err != nil {
			return fmt.Errorf("failed to put state: %w", err)
		}
	}
	// Save the list of registered options
	registeredOptionsJSON, err := json.Marshal(optionIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal registered options: %w", err)
	}
	err = ctx.GetStub().PutState("registeredOptions", registeredOptionsJSON)
	if err != nil {
		return fmt.Errorf("failed to put state for registered options: %w", err)
	}
	return nil
}

// RegisteredOptions returns the list of registered option IDs from the ledger
func (s *SmartContract) RegisteredOptions(ctx contractapi.TransactionContextInterface) ([]string, error) {
	registeredOptionsJSON, err := ctx.GetStub().GetState("registeredOptions")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registered options: %v", err)
	}

	var registeredOptions []string
	err = json.Unmarshal(registeredOptionsJSON, &registeredOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal registeredCandidates: %v", err)
	}

	return registeredOptions, nil
}



func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, voterID string, optionIDs []string) error {
	// Read the election configuration
	electionConfigBytes, err := ctx.GetStub().GetState("electionConfig")
	if err != nil {
		return fmt.Errorf("failed to read election configuration from the ledger: %v", err)
	}

	var electionConfig ElectionConfig
	err = json.Unmarshal(electionConfigBytes, &electionConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal election configuration: %v", err)
	}

	//Check if it is single option but option sent are more than one
	if electionConfig.IsSingle && len(optionIDs) > 1 {
		return fmt.Errorf("voting is single choice only, one option can be selected")
	}

	//Check if it is not abstainable but no options are sent
	if len(optionIDs) == 0 && !electionConfig.IsAbstainable {
		return fmt.Errorf("not abstainable, at least one option must be selected")
	}

	//Get the voter (fails if voter does not exist)

	voterBytes, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return fmt.Errorf("failed to read voter from the ledger: %v", err)
	}
	if voterBytes == nil {
		return fmt.Errorf("voter with ID %s does not exist", voterID)
	}

	//Get the Ballot of the voter(fails if ballot does not exist)
	ballotBytes, err := ctx.GetStub().GetState(voterID + "ballot")
	if err != nil {
		return fmt.Errorf("failed to read ballot state from the ledger: %v", err)
	}
	if ballotBytes == nil {
		return fmt.Errorf("ballot for voter with ID %s does not exist", voterID)
	}

	var ballot Ballot
	err = json.Unmarshal(ballotBytes, &ballot)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ballot: %v", err)
	}

	//check if voter already voted
	if ballot.HasVoted {
		return fmt.Errorf("voter with ID %s has already cast their vote", voterID)
	}

	// Check if each option exists in the ledger before casting the vote
	for _, optionID := range optionIDs {
		optionBytes, err := ctx.GetStub().GetState(optionID)
		if err != nil {
			return fmt.Errorf("failed to read option state from the ledger: %v", err)
		}
		if optionBytes == nil {
			return fmt.Errorf("option with ID %s does not exist", optionID)
		}
	}

	//update ballot
	ballot.HasVoted = true

	//update options voted for in the ballot if it is public(if not anonymous)
	if !electionConfig.IsAnonymous {
		ballot.OptionIDs = optionIDs
	}

	//update the state
	ballotJSON, err := json.Marshal(ballot)
	if err != nil {
		return fmt.Errorf("failed to marshal ballot json:%v", err)
	}
	err = ctx.GetStub().PutState(voterID+"ballot", ballotJSON)
	if err != nil {
		return fmt.Errorf("failed to update the ballot state: %v", err)
	}

	return nil
}


// GetVoteCount retrieves the number of votes for a given option.
func (s *SmartContract) GetVoteCount(ctx contractapi.TransactionContextInterface, optionID string) (int, error) {
	// Check if the option with the given ID exists in the ledger
	optionBytes, err := ctx.GetStub().GetState(optionID)
	if err != nil {
		return 0, fmt.Errorf("failed to read option state from the ledger: %v", err)
	}
	if optionBytes == nil {
		return 0, fmt.Errorf("option with ID %s does not exist", optionID)
	}

	// Create a composite key for the option ID to query the ledger for related ballots
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("ballot~optionID", []string{optionID})
	if err != nil {
		return 0, fmt.Errorf("failed to get iterator for option's ballots: %v", err)
	}
	defer iterator.Close()

	// Count the ballots for the given option
	voteCount := 0
	for iterator.HasNext() {
		_, err := iterator.Next()
		if err != nil {
			return 0, fmt.Errorf("failed to read ballot from the iterator: %v", err)
		}
		voteCount++
	}

	return voteCount, nil
}
