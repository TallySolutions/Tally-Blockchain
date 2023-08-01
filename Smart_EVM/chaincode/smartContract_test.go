package chaincode_test

import (
	"encoding/json"
	"fmt"
	"smart_evm/chaincode"
	"smart_evm/chaincode/mocks"

	//"strings"
	"testing"

	//"github.com/google/certificate-transparency-go/trillian/migrillian/election"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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

type ElectionConfig struct{
	IsAnonymous    bool  `json:"anonymous"`
	IsSingle       bool  `json:"single_choice"`
	IsAbstainable  bool  `json:"abstainable"`
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


	err = smartevm.InitLedger(transactionContext,true,false,true)
	require.NoError(t, err)

	// Test case 2: Error when putting the election configuration to the ledger
	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = smartevm.InitLedger(transactionContext,false,false,false)
	require.EqualError(t, err, "failed to put state: failed inserting key")

  // Test case 3: Error when fetching election configuration
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to fetch election configuration"))
  err = smartevm.InitLedger(transactionContext,true,true,false)
	require.EqualError(t,err,"failed to put state: failed inserting key")

}

type Option struct {
	ID    string `json:"id"`
	Votes int    `json:"votes"`
}


var states map[string][]byte



func GetStateStub(id string) ([]byte, error) {
   if states == nil {
        states = map[string][]byte {}
    }
    return states[id], nil
}

func AddStateStub(id string, value []byte) error {
    if states == nil {
        states = map[string][]byte {}
    }
    states[id] = value
    return nil

}
func ClearStates() {
	for k := range states {
			delete(states, k)
	}
}

type Voter struct{
	ID      string `json:"id"`
}





func TestAddVoters(t *testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub=AddStateStub
	chaincodeStub.GetStateStub=GetStateStub
	ClearStates()

	smartevm :=chaincode.SmartContract{}

	//Test Case1: Successful addition of voters
	voterIDs:=[]string{"voter1","voter2","voter3"}

	chaincodeStub.GetStateReturns(nil,nil)
	chaincodeStub.PutStateReturns(nil)

	err:=smartevm.AddVoters(transactionContext,voterIDs)
	require.NoError(t,err)

	// Test case 2: Error when registering a voter that already exists
	existingVoterID := "existingVoter"

	// Set up the mock GetState function to return a non-nil value for an existing voter
	chaincodeStub.GetStateReturns([]byte("dummyData"), nil)

	err = smartevm.AddVoters(transactionContext, []string{existingVoterID})
	require.EqualError(t, err, fmt.Sprintf("voter with ID %s already exists", existingVoterID))

  //Test Case 3: Error when writing to the ledger (PutState fails)
	voterIDsWithError := []string{"voter4","voter5"}
	chaincodeStub.GetStateReturns(nil,nil)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed to insert voter"))

	err = smartevm.AddVoters(transactionContext, voterIDsWithError)
	require.EqualError(t,err,"failed to insert voter: failed to insert voter")

	// Test case 4: Error when reading from the ledger (GetState fails)
	voterIDsWithReadError := []string{"voter6", "voter7"}

	// Set up the mock GetState function to return an error when reading voters from the ledger

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to read voter state from the ledger"))

	err = smartevm.AddVoters(transactionContext, voterIDsWithReadError)
	require.EqualError(t, err, "failed to read voter state from the ledger:failed to read voter state from the ledger")

}






func TestRegisterOptions(t *testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub=AddStateStub
	chaincodeStub.GetStateStub=GetStateStub
	ClearStates()

	smartevm :=chaincode.SmartContract{}

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
	require.EqualError(t,err,"failed to put state: failed to insert option")

	// Test case 4: Error when reading from the ledger (GetState fails)
	optionIDWithError := []string{"option6"}

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to read from world state"))

	err = smartevm.RegisterOptions(transactionContext, optionIDWithError)
	require.EqualError(t, err, "failed to read from world state: failed to read from world state")

}






func TestRegisteredOptions(t *testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub=AddStateStub
	chaincodeStub.GetStateStub=GetStateStub
	ClearStates()

	smartevm :=chaincode.SmartContract{}
	//Test Case 1: Fetch registered options successfully
	optionIDs := []string{"opt1","opt2","opt3"}

	registeredOptionsJSON ,err := json.Marshal(optionIDs)
	require.NoError(t,err)

	chaincodeStub.GetStateReturns(registeredOptionsJSON,nil)
	registeredOptions,err := smartevm.RegisteredOptions(transactionContext)
	require.NoError(t, err)
	//check if the registered option match the expected optionsIDs
	require.ElementsMatch(t,optionIDs,registeredOptions)


  //Test Case 2: failed to fetch registered options
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed to fetch registered options"))
	_, err = smartevm.RegisteredOptions(transactionContext)
	require.EqualError(t, err, "failed to fetch registered options: failed to fetch registered options")


}


func TestCastVote(t *testing.T){
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateStub=AddStateStub
	chaincodeStub.GetStateStub=GetStateStub
	ClearStates()

	smartevm :=chaincode.SmartContract{}

	//Test case 1: successful vote casting
	voterID:="voter1"
	optionIDs:=[]string{"option1","option2","option3"}

	// Set up the mock GetState function to return nil for all options (not found in the ledger)
	chaincodeStub.GetStateReturns(nil,nil)

	// Set up the mock PutState function to return nil when storing ballot and voter state
	chaincodeStub.PutStateReturns(nil)

	//set up the election configuration to be anonymous,single choice and non abstainable
	electionConfig := ElectionConfig{
		IsAnonymous: true,
		IsSingle: false,
		IsAbstainable: false,
	}
	electionConfigBytes,err :=json.Marshal(electionConfig)

	require.NoError(t,err)
	chaincodeStub.GetStateReturns(electionConfigBytes,nil)

	err=smartevm.CastVote(transactionContext,voterID,optionIDs)
	require.NoError(t,err)

	//Voting is single choice but multiple options are selected
	electionConfig = ElectionConfig{
		IsAnonymous: true,
		IsSingle: true,
		IsAbstainable: false,
	}

	electionConfigBytes,err =json.Marshal(electionConfig)

	require.NoError(t,err)
	chaincodeStub.GetStateReturns(electionConfigBytes,nil)
	optionIDs = []string{"option1", "option2"}
	err=smartevm.CastVote(transactionContext,voterID,optionIDs)
	require.EqualError(t,err,"voting is single choice only, one option can be selected")

	// Test case 3: Voting is multi-choice but no option is selected
	// Set up the election configuration to be multi-choice and abstainable
	electionConfig = ElectionConfig{
		IsAnonymous:   true,
		IsSingle:      false,
		IsAbstainable: false,
	}
	electionConfigBytes, err = json.Marshal(electionConfig)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(electionConfigBytes, nil)

	optionIDs = []string{}
	err = smartevm.CastVote(transactionContext, voterID, optionIDs)
	require.EqualError(t, err, "voting is multi choice and not abstainable,atleast one option must be selected")




}