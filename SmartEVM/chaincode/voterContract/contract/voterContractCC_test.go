package contract_test

import (
	"strings"
	"testing"

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
	"tallysolutions.com/SmartEVM/chaincode/voterContract/contract"
	"tallysolutions.com/SmartEVM/chaincode/voterContract/contract/mocks"
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

func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	voterContract := contract.SmartContract{}
	err := voterContract.InitLedger(transactionContext, true, true, true)
	require.NoError(t, err)

	//Again init should faile
	err = voterContract.InitLedger(transactionContext, true, true, true)
	require.ErrorContains(t, err, contract.ErrCCAlreadyInitialized.Error())
}

func TestInitLedgerStateFailure(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	voterContract := contract.SmartContract{}

	chaincodeStub.PutStateReturns(contract.ErrSettingState)
	err := voterContract.InitLedger(transactionContext, true, true, true)
	require.EqualError(t, err, contract.ErrSettingState.Error())
}

func TestAddVotableOption(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	//Running without initialization
	voterContract := contract.SmartContract{}
	err := voterContract.AddVotableOption(transactionContext, "Option1")
	require.Error(t, err) //Expect error

	//Initialize
	err = voterContract.InitLedger(transactionContext, true, true, true)
	require.NoError(t, err) //Expect No error

	//Now create option
	err = voterContract.AddVotableOption(transactionContext, "Option1")
	require.NoError(t, err) //Expect no error

	//create same option
	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = voterContract.AddVotableOption(transactionContext, "Option1")
	require.Error(t, err) //Expect error

	//Test state failure
	chaincodeStub.GetStateReturns(nil, contract.ErrRetrivingState)
	err = voterContract.AddVotableOption(transactionContext, "Option1")
	require.EqualError(t, err, contract.ErrRetrivingState.Error())
}

var states map[string][]byte

type PrivateData map[string][]byte

var collections map[string]PrivateData

func GetStateStub(id string) ([]byte, error) {
	if strings.HasSuffix(id, "_get_state_error_") {
		return nil, contract.ErrRetrivingState
	}
	return states[id], nil
}
func AddStateStub(id string, value []byte) error {
	if strings.HasSuffix(id, "_set_state_error_") {
		return contract.ErrSettingState
	}
	states[id] = value

	return nil
}

func PutPrivateDataStub(collection string, key string, value []byte) error {
	if collection == "_set_state_error_" {
		return contract.ErrSettingPrivate
	}
	data := PrivateData{}
	data[key] = value
	collections[collection] = data

	return nil
}

func GetPrivateDataStub(collection string, key string) ([]byte, error) {
	if collection == "_get_state_error_" {
		return nil, contract.ErrRetrivingPrivate
	}

	return collections[collection][key], nil
}

var currentIndex int

func HasNextStub() bool {
	return len(states) > currentIndex
}

var iter []queryresult.KV

func NextStub() (*queryresult.KV, error) {
	if currentIndex < len(states) {
		result := &iter[currentIndex]
		currentIndex++
		return result, nil
	}
	return nil, contract.ErrNoStateExists
}

func CloseStub() error {
	currentIndex = 0
	iter = nil
	return nil
}

func GetStateByRangeStub(arg1 string, arg2 string) (shim.StateQueryIteratorInterface, error) {
	currentIndex = 0

	iter = []queryresult.KV{}
	for id, option := range states {
		var result queryresult.KV
		result.Key = id
		result.Value = option
		iter = append(iter, result)
	}
	var iterator mocks.StateQueryIterator

	iterator.HasNextStub = HasNextStub
	iterator.NextStub = NextStub
	iterator.CloseStub = CloseStub

	return shim.StateQueryIteratorInterface(&iterator), nil
}

func SetupVote(t *testing.T, options []string, voters []string, anonymous bool, singlechoice bool, abstainable bool) (*mocks.TransactionContext, contract.SmartContract, map[string]string, error) {

	pubkeys := map[string]string{}

	states = map[string][]byte{}
	collections = map[string]PrivateData{}

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	voterContract := contract.SmartContract{}

	chaincodeStub.GetStateStub = GetStateStub
	chaincodeStub.PutStateStub = AddStateStub

	chaincodeStub.GetPrivateDataStub = GetPrivateDataStub
	chaincodeStub.PutPrivateDataStub = PutPrivateDataStub

	chaincodeStub.GetStateByRangeStub = GetStateByRangeStub

	//Initialize
	err := voterContract.InitLedger(transactionContext, anonymous, singlechoice, abstainable)
	require.NoError(t, err)

	//Set options
	for _, optionId := range options {
		err = voterContract.AddVotableOption(transactionContext, optionId)
		require.NoError(t, err)
	}

	//Set voters
	for _, voterId := range voters {
		err = voterContract.AddVoter(transactionContext, voterId)
		require.NoError(t, err)

		pubkey := AuthUser(t, voterContract, transactionContext, voterId)
		pubkeys[voterId] = pubkey
	}

	return transactionContext, voterContract, pubkeys, err
}

func StateErrorsTests(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract) {
	//Add option - set state error
	err := voterContract.AddVotableOption(transactionContext, "_set_state_error_")
	require.EqualError(t, err, contract.ErrSettingState.Error())

	//Add voter - set state error
	err = voterContract.AddVoter(transactionContext, "_set_state_error_")
	require.EqualError(t, err, contract.ErrSettingState.Error())

	//Add option - get state error
	err = voterContract.AddVotableOption(transactionContext, "_get_state_error_")
	require.EqualError(t, err, contract.ErrRetrivingState.Error())

	//Add voter - get state error
	err = voterContract.AddVoter(transactionContext, "_get_state_error_")
	require.EqualError(t, err, contract.ErrRetrivingState.Error())
}

func CastVoteTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]string, voterId string, optionId string) {

	err := voterContract.CastVote(transactionContext, voterId, keys[voterId], []string{optionId})
	require.NoError(t, err)
}

func CastVoteTestMultiChoice(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]string, voterId string, optionId1 string, optionId2 string) {

	err := voterContract.CastVote(transactionContext, voterId, keys[voterId], []string{optionId1, optionId2})
	if voterContract.SingleChoice {
		require.ErrorContains(t, err, contract.ErrNoVoteIsMoreThanOne.Error()) //Expect error
	} else {
		require.NoError(t, err)
	}
}

func AbstainVoteTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]string, voterId string) {

	err := voterContract.CastVote(transactionContext, voterId, keys[voterId], []string{})
	if voterContract.Abstainable {
		require.NoError(t, err)
	} else {
		require.ErrorContains(t, err, contract.ErrNoVoteIsZero.Error()) //Expect error
	}

}

func CastVoteWithoutAuthTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, optionId string) {

	err := voterContract.AddVoter(transactionContext, "User4")

	err = voterContract.CastVote(transactionContext, "User4", "", []string{optionId})
	require.ErrorContains(t, err, contract.ErrNotAuthorized.Error()) //Authorization Error

}

func CastVoteTwiceTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]string, voterId string, optionId string) {

	err := voterContract.CastVote(transactionContext, voterId, keys[voterId], []string{optionId})
	require.NoError(t, err)

	err = voterContract.CastVote(transactionContext, voterId, keys[voterId], []string{optionId})
	require.ErrorContains(t, err, contract.ErrAlreadyVoted.Error()) //Revote not allowed

}

func CastVoteWrongKey(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, voterId string, wrong_pubkey string, optionId string) {

	err := voterContract.CastVote(transactionContext, voterId, wrong_pubkey, []string{optionId})
	require.ErrorContains(t, err, contract.ErrNotAuthorized.Error()) //Revote not allowed

}

func ExpectedVotingResultTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, results map[string]int) {
	options, err := voterContract.GetAllOptions(transactionContext)
	require.NoError(t, err)

	for i, option := range options {
		expected := results[option.VotableId]
		require.Equal(t, expected, option.Count, "The vote count for option '%d:%s' must be %d", i, option.VotableId, expected)
	}
}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}

	return false
}
func ExpectedBallotDetailsTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, voted []string) {
	ballots, err := voterContract.GetAllBallots(transactionContext)
	require.NoError(t, err)

	for i, ballot := range ballots {
		if contains(voted, ballot.VoterId) {
			require.Equal(t, true, ballot.Casted, "The casted flag for ballot '%d:%s' must be true", i, ballot.VoterId)
			if voterContract.IsAnonymous {
				require.Equal(t, 0, len(ballot.Picks), "The picks for ballot '%d:%s' must be nil (anonymous voting)", i, ballot.VoterId)
			} else {
				require.NotEqual(t, 0, len(ballot.Picks), "The picks for ballot '%d:%s' must not be nil (public voting)", i, ballot.VoterId)
			}
			require.NotEqual(t, int64(0), ballot.Timestamp, "The timestamp for ballot '%d:%s' can not be zero", i, ballot.VoterId)
		} else {
			require.Equal(t, false, ballot.Casted, "The casted flag for ballot '%d:%s' must be true", i, ballot.VoterId)
			require.Equal(t, int64(0), ballot.Timestamp, "The timestamp for ballot '%d:%s' must be zero", i, ballot.VoterId)
		}

	}
}

// Test Anonymous, Single-Choice and Abstainable Election
func TestVote_Anonymous_SingleChoice_Abstainable(t *testing.T) {

	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, true, true, true)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTwiceTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	CastVoteWithoutAuthTest(t, transactionContext, voterContract, "Option1") //Will not register

	CastVoteWrongKey(t, transactionContext, voterContract, "User2", keys["User1"], "Option2") //Will not refgister

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User2")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User3", "Option1", "Option2") //Will not register

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		contract.Abstained: 1,
		"Option1":          1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
		"User2",
	})

}

// Test Anonymous, Single-Choice and Not-Abstainable Election
func TestVote_Anonymous_SingleChoice_NonAbstainable(t *testing.T) {
	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, true, true, false)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User2", "Option2", "Option3") //Will not register

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User2") //Will not register

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		"Option1": 1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
	})

}

// Test Anonymous, Multi-Choice and Not-Abstainable Election
func TestVote_Anonymous_MultiChoice_Abstainable(t *testing.T) {
	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, true, false, true)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User2", "Option2", "Option3")

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User3")

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		contract.Abstained: 1,
		"Option1":          1,
		"Option2":          1,
		"Option3":          1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
		"User2",
		"User3",
	})

}

// Test Anonymous, Multi-Choice and Not-Abstainable Election
func TestVote_Anonymous_MultiChoice_NonAbstainable(t *testing.T) {
	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, true, false, false)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User2", "Option2", "Option3")

	CastVoteTest(t, transactionContext, voterContract, keys, "User3", "Option2")

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User2")

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		"Option1": 1,
		"Option2": 2,
		"Option3": 1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
		"User2",
		"User3",
	})

}

// Test Public, Single-Choice and Abstainable Election
func TestVote_Public_SingleChoice_Abstainable(t *testing.T) {
	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, false, true, true)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User2")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User3", "Option1", "Option2") //Will not register

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		contract.Abstained: 1,
		"Option1":          1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
		"User2",
	})
}

// Test Public, Single-Choice and Not-Abstainable Election
func TestVote_Public_SingleChoice_NonAbstainable(t *testing.T) {
	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, false, true, false)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User2")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User3", "Option1", "Option2") //Will fail to register

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		"Option1": 1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
	})

}

// Test Anonymous, Multi-Choice and Not-Abstainable Election
func TestVote_Public_MultiChoice_Abstainable(t *testing.T) {
	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, false, false, true)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User2")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User3", "Option1", "Option3")

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		contract.Abstained: 1,
		"Option1":          2,
		"Option3":          1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
		"User2",
		"User3",
	})

}

// Test Public, Multi-Choice and Not-Abstainable Election
func TestVote_Public_MultiChoice_NonAbstainable(t *testing.T) {
	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, false, false, false)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User2")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User3", "Option1", "Option2")

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		"Option1": 2,
		"Option2": 1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
		"User3",
	})

}

func AuthUser(t *testing.T, voterContract contract.SmartContract, transactionContext *mocks.TransactionContext, userId string) string {

	//First authenticate user

	//Generate a key pair
	// provision key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_bytes := x509.MarshalPKCS1PublicKey(&publicKey)
	publicKey_base64 := base64.StdEncoding.EncodeToString(publicKey_bytes)

	//Sign
	msgHash := sha512.New()
	_, err = msgHash.Write([]byte(userId))
	require.NoError(t, err) //Expect No error

	msgHashSum := msgHash.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA512, msgHashSum, nil)
	require.NoError(t, err) //Expect No error

	signature_base64 := base64.StdEncoding.EncodeToString(signature)

	pub, err := voterContract.AuthVoter(transactionContext, "_get_state_error_", publicKey_base64, signature_base64)
	require.ErrorContains(t, err, contract.ErrRetrivingState.Error()) //Expect error

	pub, err = voterContract.AuthVoter(transactionContext, userId, publicKey_base64, signature_base64)
	require.NoError(t, err) //Expect No error

	return pub

}
