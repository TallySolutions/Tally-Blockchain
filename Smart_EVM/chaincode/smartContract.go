package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	//"github.com/hyperledger/fabric-sdk-go/pkg/common/options"
)

type SmartContract struct {
    contractapi.Contract
}

type ElectionConfig struct{
	IsAnonymous    bool  `json:"anonymous"`
	IsSingle       bool  `json:"single_choice"`
	IsAbstainable  bool  `json:"abstainable"`
}

type Option struct{
	ID      string `json:"id"`
	Votes   int    `json:"votes"`
}

type Voter struct{
	ID      string `json:"id"`
}

type Ballot struct{
	VoterID    string      `json:"voterID"`
	OptionID   string      `json:"optionID"`
	HasVoted   bool        `json:"hasVoted"`
	Timestamp  time.Time   `json:"timestamp"`
}


func(s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, isAnonymous bool, isSingle bool, isAbstainable bool) error{
	  //create the election configuration
    electionConfig := ElectionConfig{
			IsAnonymous     :  isAnonymous,
			IsSingle        :  isSingle,
			IsAbstainable   :  isAbstainable,
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
			{ID:"1" , Votes:0},
			{ID:"2" , Votes:0},
			{ID:"3" , Votes:0},
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
  for _, voterID:=range voterIDs{
		existingVoterJSON ,err :=ctx.GetStub().GetState(voterID)
		if err!=nil{
			return fmt.Errorf("failed to read voter state from the ledger:%w",err)
		}
		if existingVoterJSON!=nil{
			return fmt.Errorf("voter with ID %s already exists", voterID)
		}

		voter := Voter{
			ID:       voterID,
		}

		voterJSON,err:=json.Marshal(voter)
		if err!=nil{
			return fmt.Errorf("failed to marshal voter JSON :%w",err)
		}

		err=ctx.GetStub().PutState(voterID,voterJSON)
		if err != nil {
			return fmt.Errorf("failed to insert voter: %v", err)
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




func (s* SmartContract) CastVote(ctx contractapi.TransactionContextInterface , voterID string , optionIDs []string) error{
	//check if election is anonymous
	electionConfigBytes, err:=ctx.GetStub().GetState("electionConfig")
	if err!=nil{
		return fmt.Errorf("failed to read election configuration from the ledger: %v",err)
	}

	var electionConfig ElectionConfig
	err =json.Unmarshal(electionConfigBytes,&electionConfig)
	if err!=nil{
		return fmt.Errorf("failed to unmarshal election configuration: %v",err)
	}

	isSingleChoice:=electionConfig.IsSingle

	if isSingleChoice && len(optionIDs) >1{
		return fmt.Errorf("voting is single choice only, one option can be selected")
	}

	if !isSingleChoice && len(optionIDs)==0 && !electionConfig.IsAbstainable{
		return fmt.Errorf("voting is multi choice and not abstainable,atleast one option must be selected")
	}

	//check if the voter with the given ID already exists in the ledger
	voterBytes, err:=ctx.GetStub().GetState(voterID)
	if err!=nil{
		return fmt.Errorf("failed to read voter from the ledger: %v",err)
	}
	if !electionConfig.IsAnonymous && voterBytes!=nil{
		return fmt.Errorf("voter with ID %s already exists",voterID)
	}

	//if voting is abstaible and no options are slected , add "None of the Above" option
	if electionConfig.IsAbstainable && len(optionIDs)==0{
		optionIDs = append(optionIDs, "None of the Above")
	}

	//create a ballot for each selected option and store it in the ledger
	for _, optionID :=range optionIDs{
		//check if the option with the given ID exist in the ledger
		optionBytes,err :=ctx.GetStub().GetState(optionID)
		if err!=nil{
			return fmt.Errorf("failed to read option state from the ledger: %v",err)
		}
		if optionBytes==nil{
			return fmt.Errorf("option with ID %s does not exist",optionID)
		}

		//create a ballot with voter's ID , option Id , hasVoted set to true and timestamp
		ballot:=Ballot{
			VoterID: voterID,
			OptionID: optionID,
			HasVoted: true,
			Timestamp: time.Now(),
		}

		//marshal the ballot to json and store it to the ledger
		ballotJSON,err:=json.Marshal(ballot)
		if err!=nil{
			return fmt.Errorf("failed to marshal ballot JSON: %v",err)
		}
		err=ctx.GetStub().PutState(voterID+optionID , ballotJSON)
		if err!=nil{
			return fmt.Errorf("failed to put ballot state:%v",err)
		}
	}

	//Mark Voter as voted (if voting is not anonymous)
	if !electionConfig.IsAnonymous{
		voter:=Voter{
			ID:voterID,
		}

		voterJSON,err := json.Marshal(voter)
		if err!=nil{
			return fmt.Errorf("failed to marshal voter JSON:%v",err)
		}
		err= ctx.GetStub().PutState(voterID,voterJSON)
		if err!=nil{
			return fmt.Errorf("failed to put voter state: %v",err)
		}
	}

	return nil
}

