package chaincode_test

import (
	"encoding/json"
	"fmt"
	"smart_evm/chaincode"
	"smart_evm/chaincode/mocks"
	"strings"
	"testing"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

type ElectionConfig struct {
	IsAnonymous   bool `json:"anonymous"`
	IsSingle      bool `json:"single_choice"`
	IsAbstainable bool `json:"abstainable"`
}



func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	smartevm := chaincode.SmartContract{}

	// Test case 1: Successful initialization of the ledger
	isAnonymous := true
	isSingle := true
	isAbstainable := true

	// Set up the mock GetState function to return the election configuration data
	electionConfig := ElectionConfig{
		IsAnonymous:   isAnonymous,
		IsSingle:      isSingle,
		IsAbstainable: isAbstainable,
	}

	electionConfigBytes, err := json.Marshal(electionConfig)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(electionConfigBytes, nil)

	// Call the InitLedger function
	err = smartevm.InitLedger(transactionContext, isAnonymous, isSingle, isAbstainable)
	require.NoError(t, err)

	// Verify that the election configuration is saved to the ledger
	electionConfigBytes, err = chaincodeStub.GetState("electionConfig")
	require.NoError(t, err)
	require.NotNil(t, electionConfigBytes)

	var actualElectionConfig ElectionConfig
	err = json.Unmarshal(electionConfigBytes, &actualElectionConfig)
	require.NoError(t, err)

	require.Equal(t, electionConfig, actualElectionConfig)

	err = smartevm.InitLedger(transactionContext, true, false, true)
	require.NoError(t, err)

	// Test case 2: Error when putting the election configuration to the ledger
	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = smartevm.InitLedger(transactionContext, false, false, false)
	require.EqualError(t, err, "failed to put state: failed inserting key")

	// Test case 3: Error when fetching election configuration
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to fetch election configuration"))
	err = smartevm.InitLedger(transactionContext, true, true, false)
	require.EqualError(t, err, "failed to put state: failed inserting key")

}


type Option struct {
	ID    string `json:"id"`
	Votes int    `json:"votes"`
}


var states map[string][]byte

func GetStateStub(id string) ([]byte, error) {
	if states == nil {
		states = map[string][]byte{}
	}
	return states[id], nil
}


func AddStateStub(id string, value []byte) error {
	if states == nil {
		states = map[string][]byte{}
	}
	states[id] = value
	return nil

}
func ClearStates() {
	for k := range states {
		delete(states, k)
	}
}


type Voter struct {
	ID string `json:"id"`
}



func TestAddVoters(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	//Test Case1: Successful addition of voters
	voterIDs := []string{"voter1", "voter2", "voter3"}

	chaincodeStub.GetStateReturns(nil, nil)
	chaincodeStub.PutStateReturns(nil)

	err := smartevm.AddVoters(transactionContext, voterIDs)
	require.NoError(t, err)

	// Test case 2: Error when registering a voter that already exists
	existingVoterID := "existingVoter"

	// Set up the mock GetState function to return a non-nil value for an existing voter
	chaincodeStub.GetStateReturns([]byte("dummyData"), nil)

	err = smartevm.AddVoters(transactionContext, []string{existingVoterID})
	require.EqualError(t, err, fmt.Sprintf("voter with ID %s already exists", existingVoterID))

	//Test Case 3: Error when writing to the ledger (PutState fails)
	voterIDsWithError := []string{"voter4", "voter5"}
	chaincodeStub.GetStateReturns(nil, nil)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed to insert voter"))

	err = smartevm.AddVoters(transactionContext, voterIDsWithError)
	require.EqualError(t, err, "failed to insert voter: failed to insert voter")

	// Test case 4: Error when reading from the ledger (GetState fails)
	voterIDsWithReadError := []string{"voter6", "voter7"}

	// Set up the mock GetState function to return an error when reading voters from the ledger

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to read voter state from the ledger"))

	err = smartevm.AddVoters(transactionContext, voterIDsWithReadError)
	require.EqualError(t, err, "failed to read voter state from the ledger:failed to read voter state from the ledger")

}



func TestRegisterOptions(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	// Test case 1: Register multiple options successfully
	optionIDs := []string{"option1", "option2", "option3"}

	err := smartevm.RegisterOptions(transactionContext, optionIDs)
	require.NoError(t, err)

	// Verify that each option is saved to the ledger with an initial vote count of 0
	for _, optionID := range optionIDs {
		optionBytes, err := chaincodeStub.GetState(optionID)
		require.NoError(t, err)
		require.NotNil(t, optionBytes)

		var option Option
		err = json.Unmarshal(optionBytes, &option)
		require.NoError(t, err)

		require.Equal(t, 0, option.Votes)
	}

	// Verify that the list of registered options is saved to the ledger
	registeredOptionsBytes, err := chaincodeStub.GetState("registeredOptions")
	require.NoError(t, err)
	require.NotNil(t, registeredOptionsBytes)

	var registeredOptions []string
	err = json.Unmarshal(registeredOptionsBytes, &registeredOptions)
	require.NoError(t, err)

	require.ElementsMatch(t, optionIDs, registeredOptions)

	// Test case 2: Error when registering an option that already exists
	existingOptionID := "existingOption"

	// Set up the mock GetState function to return a non-nil value for an existing option
	chaincodeStub.GetStateReturns([]byte("dummyData"), nil)

	err = smartevm.RegisterOptions(transactionContext, []string{existingOptionID})
	require.EqualError(t, err, fmt.Sprintf("option with ID %s already exists", existingOptionID))

	// Test case 3: Error when writing to the ledger (PutState fails)
	optionIDsWithError := []string{"option4", "option5"}

	// Set up the mock GetState function to return nil for options that don't exist in the ledger
	chaincodeStub.GetStateReturns(nil, nil)

	// Set up the mock PutState function to return an error when registering options

	chaincodeStub.PutStateReturns(fmt.Errorf("failed to insert option"))

	err = smartevm.RegisterOptions(transactionContext, optionIDsWithError)
	require.EqualError(t, err, "failed to put state: failed to insert option")

	// Test case 4: Error when reading from the ledger (GetState fails)
	optionIDWithError := []string{"option6"}

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to read from world state"))

	err = smartevm.RegisterOptions(transactionContext, optionIDWithError)
	require.EqualError(t, err, "failed to read from world state: failed to read from world state")

}



func TestRegisteredOptions(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}
	//Test Case 1: Fetch registered options successfully
	optionIDs := []string{"opt1", "opt2", "opt3"}

	registeredOptionsJSON, err := json.Marshal(optionIDs)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(registeredOptionsJSON, nil)
	registeredOptions, err := smartevm.RegisteredOptions(transactionContext)
	require.NoError(t, err)
	//check if the registered option match the expected optionsIDs
	require.ElementsMatch(t, optionIDs, registeredOptions)

	//Test Case 2: failed to fetch registered options
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to fetch registered options"))
	_, err = smartevm.RegisteredOptions(transactionContext)
	require.EqualError(t, err, "failed to fetch registered options: failed to fetch registered options")

}


type Ballot struct {
	VoterID   string    `json:"voterID"`
	OptionIDs []string  `json:"optionIDs"`
	HasVoted  bool      `json:"hasVoted"`
	Timestamp time.Time `json:"timestamp"`
}


func TestCastVotefalsefalsefalse(t* testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	//Testcase: false false false(non-anonymous , non-single,non-abstainable)
	voterID:="voter1"
	optionIDs:=[]string{"option1","option2"}

	//set up the election configuration
	err:=smartevm.InitLedger(transactionContext,false,false,false)
	require.NoError(t,err)

	//add voter
	err=smartevm.AddVoters(transactionContext,[]string{voterID})
	require.NoError(t,err)

	//register options
	err=smartevm.RegisterOptions(transactionContext,optionIDs)
	require.NoError(t,err)

	// Cast vote with a non-existing option ID
	err = smartevm.CastVote(transactionContext, voterID, []string{"noSuchOption"})
	require.EqualError(t, err, "option with ID noSuchOption does not exist")


	//Cast Vote Successfully
	err=smartevm.CastVote(transactionContext,voterID,optionIDs)
	require.NoError(t,err)

	//Expect error:voter already exist
	err=smartevm.AddVoters(transactionContext,[]string{voterID})
	require.EqualError(t,err,fmt.Sprintf("voter with ID %s already exists",voterID))

	//Expect error:option already exist
	err=smartevm.RegisterOptions(transactionContext,optionIDs)
	require.EqualError(t,err,fmt.Sprintf("option with ID %s already exists",optionIDs[0]))



	// Expect error: voter already cast their vote
	err = smartevm.CastVote(transactionContext, voterID, optionIDs)
	require.EqualError(t, err, fmt.Sprintf("voter with ID %s has already cast their vote", voterID)) 
}



func TestCastVotefalsetruetrue(t* testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	//Testcase: Successful vote casting for single choice , abstainable election(false,true,true)
	voterID:="voter1"
	voterID1:="voter2"
	optionIDs:=[]string{"option1","option2","option3"}


	// Set up the election configuration to be single-choice, abstainable, and non-anonymous
	err := smartevm.InitLedger(transactionContext, false, true, true)
	require.NoError(t, err)


	// Register options
	smartevm.RegisterOptions(transactionContext, optionIDs)
	require.NoError(t, err)

	// Register a voter
	err = smartevm.AddVoters(transactionContext, []string{voterID,voterID1})
	require.NoError(t, err)

	// Cast vote with multiple options for single-choice election
	err = smartevm.CastVote(transactionContext, voterID, optionIDs)
	require.Error(t, err, "voting is single choice only, one option can be selected")

	// Cast vote without selecting any options (abstainable voting)
	err = smartevm.CastVote(transactionContext, voterID1, []string{})
	require.NoError(t,err)

	// Cast vote with a single option for single-choice election
	err = smartevm.CastVote(transactionContext, voterID, []string{"option1"})
	require.NoError(t, err)

	// Attempt to cast vote for a voter with no ballot state
	err=smartevm.CastVote(transactionContext,"voter3",[]string{})
	require.EqualError(t,err,"voter with ID voter3 does not exist")
}



func TestCastVotefalsetruefalse(t* testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	optionsIds := []string{"option1","option2","options3"}

	//set up the election configuration to be non-anonymous,single,non-abstainable
	err:=smartevm.InitLedger(transactionContext,false,true,false)
	require.NoError(t,err)

	err=smartevm.RegisterOptions(transactionContext,optionsIds)
	require.NoError(t,err)

	//register a voter
	voterID:="voter1"
	err=smartevm.AddVoters(transactionContext,[]string{voterID})
	require.NoError(t,err)

	//Cast with multiple options for single choice voting
	err=smartevm.CastVote(transactionContext,voterID,optionsIds)
	require.Error(t,err,"voting is single choice only, one option can be selected")

	//Cast with a single option for single choice voting
	err=smartevm.CastVote(transactionContext,voterID,[]string{"option1"})
	require.NoError(t,err)

	//Cast vote without selecting any options(non-abstainable)
	err=smartevm.CastVote(transactionContext,voterID,[]string{})
	require.Error(t,err,"not abstainable, at least one option must be selected")

}



func TestCastVotefalsefalsetrue(t* testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	//set up the election configuration to be non-anonymous,multi,abstainable
	err:=smartevm.InitLedger(transactionContext,false,false,true)
	require.NoError(t,err)

	//Register options
  optionIDs:= []string{"option1","option2","options3"}
	err=smartevm.RegisterOptions(transactionContext,optionIDs)
	require.NoError(t,err)

	//Register voter
	voterID:=[]string{"voter1","voter2","voter3"}
	err=smartevm.AddVoters(transactionContext,voterID)
	require.NoError(t,err)

	// Cast vote with multiple options for multi-choice election
	err=smartevm.CastVote(transactionContext,"voter1",optionIDs)
	require.NoError(t,err)

  //Cast vote without selecting any options(abstainable)
	err=smartevm.CastVote(transactionContext,"voter2",[]string{})
	require.NoError(t,err)

	//Cast with a single option in multi choice voting
	err=smartevm.CastVote(transactionContext,"voter3",[]string{"option1"})
	require.NoError(t,err)
}



func TestCastVotetruefalsefalse(t* testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	//set up the election configuration to be anonymous,multi,non-abstainable
	err:=smartevm.InitLedger(transactionContext,true,false,false)
	require.NoError(t,err)

	//Register options
  optionIDs:= []string{"option1","option2","option3"}
	err=smartevm.RegisterOptions(transactionContext,optionIDs)
	require.NoError(t,err)

	//Register voter
	voterID:=[]string{"voter1","voter2","voter3"}
	err=smartevm.AddVoters(transactionContext,voterID)
	require.NoError(t,err)

	// Cast vote with multiple options for multi-choice election
	err=smartevm.CastVote(transactionContext,"voter1",optionIDs)
	require.NoError(t,err)

	//Cast vote without selecting any options(non-abstainable)
	err=smartevm.CastVote(transactionContext,"voter2",[]string{})
	require.Error(t,err,"not abstainable, at least one option must be selected")

	//Cast for single choice in multi choice voting
	err = smartevm.CastVote(transactionContext, "voter2", []string{"option1"})
	require.NoError(t, err)

	// Cast vote with a different set of options
	err = smartevm.CastVote(transactionContext, "voter3", []string{"option2", "option3"})
	require.NoError(t, err)

}



func TestCastVotetruefalsetrue(t* testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	optionIDs := []string{"option1", "option2", "option3"}

	// Set up the election configuration to be multi-choice, abstainable, and anonymous
	err := smartevm.InitLedger(transactionContext, true, false, true)
	require.NoError(t, err)

	// Register options
	smartevm.RegisterOptions(transactionContext, optionIDs)
	require.NoError(t, err)

	// Add voters
	voterIDs := []string{"voter1", "voter2", "voter3"}
	err = smartevm.AddVoters(transactionContext, voterIDs)
	require.NoError(t, err)

	// Cast vote with multiple options for multi-choice election
	err = smartevm.CastVote(transactionContext, "voter1", optionIDs)
	require.NoError(t, err)

	// Cast vote without selecting any options (abstainable voting)
	err = smartevm.CastVote(transactionContext, "voter2", []string{})
	require.NoError(t, err)

	// Cast vote with a single option in multi-choice voting
	err = smartevm.CastVote(transactionContext, "voter3", []string{"option1"})
	require.NoError(t, err)
}



func TestCastVotetruetruefalse(t *testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	optionIDs := []string{"option1", "option2", "option3"}

	// Set up the election configuration to be single-choice, non-abstainable, and anonymous
	err := smartevm.InitLedger(transactionContext, true, true, false)
	require.NoError(t, err)

	// Register options
	smartevm.RegisterOptions(transactionContext, optionIDs)
	require.NoError(t, err)

	// Add voters
	voterIDs := []string{"voter1", "voter2", "voter3"}
	err = smartevm.AddVoters(transactionContext, voterIDs)
	require.NoError(t, err)

	//multi choice voting in single choice type
	err=smartevm.CastVote(transactionContext,"voter1",optionIDs)
	require.Error(t,err,"voting is single choice only, one option can be selected")

	//no option selected in non-abtainable voting
	err=smartevm.CastVote(transactionContext,"voter2",[]string{})
	require.Error(t,err,"not abstainable, at least one option must be selected")

	//single option selected in single choice voting
	err=smartevm.CastVote(transactionContext,"voter3",[]string{"option1"})
	require.NoError(t,err)

}



func TestCastVotetruetruetrue(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	// Set up the election configuration to be anonymous , single-choice and abstainable
	err := smartevm.InitLedger(transactionContext, true, true, true)
	require.NoError(t, err)

	//register options
	optionIDs := []string{"option1", "option2", "option3"}
	smartevm.RegisterOptions(transactionContext, optionIDs)
	require.NoError(t, err)

	//register voter
	voterIDs := []string{"voter1", "voter2", "voter3"}
	err = smartevm.AddVoters(transactionContext, voterIDs)
	require.NoError(t, err)

	// Cast vote with multiple options for single-choice election
	err = smartevm.CastVote(transactionContext, "voter1", optionIDs)
	require.Error(t, err, "voting is single choice only, one option can be selected")

	// Cast vote without selecting any options (abstainable voting)
	err = smartevm.CastVote(transactionContext, "voter2", []string{})
	require.NoError(t, err)

	// Cast vote with a single option for single-choice election
	err = smartevm.CastVote(transactionContext, "voter3", []string{"option1"})
	require.NoError(t, err)

}

func TestCastVoteFailedToUpdateBallotState(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub

	// Clear any existing states
	ClearStates()

	smartevm := chaincode.SmartContract{}

	// Set up the election configuration
	err := smartevm.InitLedger(transactionContext, false, false, false)
	require.NoError(t, err)

	// Add a voter
	voterID := "voter1"
	err = smartevm.AddVoters(transactionContext, []string{voterID})
	require.NoError(t, err)

	// Register an option
	optionID := "option1"
	err = smartevm.RegisterOptions(transactionContext, []string{optionID})
	require.NoError(t, err)

	// Mock the error while updating the ballot state
	chaincodeStub.PutStateStub = func(key string, value []byte) error {
		return fmt.Errorf("failed to update ballot state")
	}

	// Attempt to cast a vote
	err = smartevm.CastVote(transactionContext, voterID, []string{optionID})
	require.EqualError(t, err, "failed to update the ballot state: failed to update ballot state")

	// Clear any remaining states
	ClearStates()
}






func TestGetVoteCount(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub = AddStateStub
	chaincodeStub.GetStateStub = GetStateStub
	ClearStates()

	smartevm := chaincode.SmartContract{}

	// Test case 1: Successful retrieval of vote count for an existing option
	optionID := "option1"
	voteCount := 5

	// Set up the state of the ledger to show that the option exists
	option := Option{
		ID: optionID,
	}
	optionBytes, err := json.Marshal(option)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(optionBytes, nil)

	// Set up the iterator for the option's ballots
	ballotIterator := &mocks.StateQueryIterator{}
	chaincodeStub.GetStateByPartialCompositeKeyReturns(ballotIterator, nil)
	ballotIterator.HasNextReturnsOnCall(0, true)
	ballotIterator.HasNextReturnsOnCall(1, true)
	ballotIterator.HasNextReturnsOnCall(2, true)
	ballotIterator.HasNextReturnsOnCall(3, true)
	ballotIterator.HasNextReturnsOnCall(4, true)
	ballotIterator.HasNextReturnsOnCall(5, false)
	ballotIterator.NextReturnsOnCall(0, &queryresult.KV{Value: []byte("ballot1")}, nil)
	ballotIterator.NextReturnsOnCall(1, &queryresult.KV{Value: []byte("ballot2")}, nil)
	ballotIterator.NextReturnsOnCall(2, &queryresult.KV{Value: []byte("ballot3")}, nil)
	ballotIterator.NextReturnsOnCall(3, &queryresult.KV{Value: []byte("ballot4")}, nil)
	ballotIterator.NextReturnsOnCall(4, &queryresult.KV{Value: []byte("ballot5")}, nil)

	// Invoke GetVoteCount function
	count, err := smartevm.GetVoteCount(transactionContext, optionID)
	require.NoError(t, err)
	require.Equal(t, voteCount, count)

	// Test case 2: Retrieval of vote count for a non-existing option
	nonExistingOptionID := "nonExistingOption"

	// Set up the state of the ledger to show that the option does not exist
	chaincodeStub.GetStateReturns(nil, nil)

	// Invoke GetVoteCount function for non-existing option
	count, err = smartevm.GetVoteCount(transactionContext, nonExistingOptionID)
	require.EqualError(t, err, fmt.Sprintf("option with ID %s does not exist", nonExistingOptionID))
	require.Equal(t, 0, count)

	// Test case 3: Error when option state retrieval fails
	failingOptionID := "failingOption"

	// Set up the state of the ledger to return an error
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to retrieve option state"))

	// Invoke GetVoteCount function for the failing option
	_, err = smartevm.GetVoteCount(transactionContext, failingOptionID)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "failed to read option state from the ledger: failed to retrieve option state"))
}

