package contract_test

import (
	"fmt"
	"strings"
	"testing"

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	"github.com/SmartEVM/chaincode/voterContract/contract"
	"github.com/SmartEVM/chaincode/voterContract/contract/mocks"
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

type KeyPair struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

var current_test = 1
var total_test = 14 //Change if you add new test
var current_subtest = 0
var total_subtest = 0

func logTest(name string, subtests int) {
	fmt.Printf("[%d/%d] Running Test '%s' ...\n", current_test, total_test, name)
	current_test++
	total_subtest = subtests
	current_subtest = 1
}
func logSubTest(name string) {
	fmt.Printf("      [%d/%d] Running Subtest '%s' ...\n", current_subtest, total_subtest, name)
	current_subtest++
}

func TestInitLedger(t *testing.T) {
	logTest("TestInitLedger", 0)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	voterContract := contract.SmartContract{}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_base64 := EncodePublicKey(publicKey)

	publicKey_base64, err = voterContract.InitLedger(transactionContext, publicKey_base64, true, true, true, true)
	require.NoError(t, err)

	//Again init should fail
	publicKey_base64, err = voterContract.InitLedger(transactionContext, publicKey_base64, true, true, true, false)
	require.ErrorContains(t, err, contract.ErrCCAlreadyInitialized.Error())
}

func TestInitLedgerFailureToAddAbsytained(t *testing.T) {
	logTest("TestInitLedgerFailureToAddAbsytained", 0)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	voterContract := contract.SmartContract{}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_base64 := EncodePublicKey(publicKey)

	chaincodeStub.GetStateReturns(nil, contract.ErrSettingState)
	publicKey_base64, err = voterContract.InitLedger(transactionContext, publicKey_base64, true, true, true, true)
	require.Error(t, contract.ErrSettingState)

}

func TestInitLedgerStateFailure(t *testing.T) {
	logTest("TestInitLedgerStateFailure", 0)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	voterContract := contract.SmartContract{}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_base64 := EncodePublicKey(publicKey)

	chaincodeStub.PutStateReturns(contract.ErrSettingState)
	publicKey_base64, err = voterContract.InitLedger(transactionContext, publicKey_base64, true, true, true, false)
	require.EqualError(t, err, contract.ErrSettingState.Error())
}

func TestAddVotableOption(t *testing.T) {
	logTest("TestAddVotableOption", 0)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	//Running without initialization
	voterContract := contract.SmartContract{}
	err := voterContract.AddVotableOption(transactionContext, "Option1")
	require.Error(t, err) //Expect error

	//Initialize
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_base64 := EncodePublicKey(publicKey)
	publicKey_base64, err = voterContract.InitLedger(transactionContext, publicKey_base64, true, true, true, true)
	require.NoError(t, err) //Expect No error

	//Now create option
	err = voterContract.AddVotableOption(transactionContext, "Option1")
	require.NoError(t, err) //Expect no error

	//create same option
	votableOption := contract.VotableOption{
		VotableId: "Option1",
		Count:     0,
	}
	votableOptionJSON, err := json.Marshal(votableOption)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte(votableOptionJSON), nil)
	err = voterContract.AddVotableOption(transactionContext, "Option1")
	require.ErrorContains(t, err, contract.ErrVotingOptionAlreadyExists.Error()) //Expect error

	//Test state failure
	chaincodeStub.GetStateReturns(nil, contract.ErrRetrivingState)
	err = voterContract.AddVotableOption(transactionContext, "Option1")
	require.EqualError(t, err, contract.ErrRetrivingState.Error())
}

func TestAddVotersPartialFailure(t *testing.T) {
	logTest("TestAddVotersPartialFailure", 0)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	states = map[string][]byte{}

	chaincodeStub.GetStateStub = GetStateStub
	chaincodeStub.PutStateStub = AddStateStub

	voterContract := contract.SmartContract{}

	//Initialize
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_base64 := EncodePublicKey(publicKey)

	publicKey_base64, err = voterContract.InitLedger(transactionContext, publicKey_base64, true, true, true, false)
	require.NoError(t, err)

	//Set voters
	//Sign data with the private key

	voters := []string{"User1", "User2", "User3", "User2"}

	signature_bytes := SignDataBytes(t, privateKey, strings.Join(voters, ","))

	//Encryot data with public key
	signature_enc, err := contract.EncryptOAEP(DecodePublicKey(t, publicKey_base64), signature_bytes)
	require.NoError(t, err)

	//1 out of 4 should fail
	err, errors := voterContract.AddVoters(transactionContext, voters, base64.StdEncoding.EncodeToString(signature_enc))
	require.ErrorContains(t, err, contract.ErrCouldAddAddAllVoters.Error())
	require.Equal(t, 1, len(errors))                                            //Should be one error
	require.ErrorContains(t, errors[0], contract.ErrVoterAlreadyExists.Error()) //Error should be voter already exists
}

func TestAddVotersKeyErrors(t *testing.T) {
	logTest("TestAddVotersPartialFailure", 0)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	voterContract := contract.SmartContract{}

	//Initialize
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_base64 := EncodePublicKey(publicKey)

	publicKey_base64, err = voterContract.InitLedger(transactionContext, publicKey_base64, true, true, true, false)
	require.NoError(t, err)

	//Set voters
	//Sign data with the private key

	voters := []string{"User1", "User2"}

	//Sign with wrong private key test
	privateKey2, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error
	signature_bytes := SignDataBytes(t, privateKey2, strings.Join(voters, ","))

	//Encryot data with public key
	signature_enc, err := contract.EncryptOAEP(DecodePublicKey(t, publicKey_base64), signature_bytes)
	require.NoError(t, err)

	//Should fail - wrong signature
	err, _ = voterContract.AddVoters(transactionContext, voters, base64.StdEncoding.EncodeToString(signature_enc))
	require.ErrorContains(t, err, contract.ErrSignatureValidation.Error())

	//Encrypt with wrong public key test
	signature_bytes = SignDataBytes(t, privateKey, strings.Join(voters, ","))
	//Encryot data with public key
	signature_enc, err = contract.EncryptOAEP(&privateKey2.PublicKey, signature_bytes)
	require.NoError(t, err)
	//Should fail - decryption error
	err, _ = voterContract.AddVoters(transactionContext, voters, base64.StdEncoding.EncodeToString(signature_enc))
	require.ErrorContains(t, err, contract.ErrDecryption.Error())

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

func SetupVote(t *testing.T, options []string, voters []string, anonymous bool, singlechoice bool, abstainable bool) (*mocks.TransactionContext, contract.SmartContract, map[string]KeyPair, error) {
	fmt.Println("  Setting up ...")
	var bar contract.Bar
	bar.NewOption(100, "  ")
	keys := map[string]KeyPair{}

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
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_base64 := EncodePublicKey(publicKey)

	publicKey_base64, err = voterContract.InitLedger(transactionContext, publicKey_base64, anonymous, singlechoice, abstainable, false)
	require.NoError(t, err)

	//Set options
	for _, optionId := range options {
		err = voterContract.AddVotableOption(transactionContext, optionId)
		require.NoError(t, err)
	}

	bar.Play(1)

	//Set voters
	//Sign data with the private key
	signature_bytes := SignDataBytes(t, privateKey, strings.Join(voters, ","))

	//Encryot data with public key
	signature_enc, err := contract.EncryptOAEP(DecodePublicKey(t, publicKey_base64), signature_bytes)
	require.NoError(t, err)

	err, _ = voterContract.AddVoters(transactionContext, voters, base64.StdEncoding.EncodeToString(signature_enc))
	require.NoError(t, err)

	bar.Play(10)

	//Authorize Voters
	len_voters := len(voters)
	for i, voterId := range voters {
		keys[voterId] = AuthUser(t, voterContract, transactionContext, voterId)
		progress := int64((((i + 1) * 90) / len_voters) + 10)
		bar.Play(progress)
	}
	bar.Finish()
	return transactionContext, voterContract, keys, err
}

func StateErrorsTests(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract) {
	logSubTest("StateErrorsTests")

	//Add option - set state error
	err := voterContract.AddVotableOption(transactionContext, "_set_state_error_")
	require.EqualError(t, err, contract.ErrSettingState.Error())

	//Add option - get state error
	err = voterContract.AddVotableOption(transactionContext, "_get_state_error_")
	require.EqualError(t, err, contract.ErrRetrivingState.Error())

}

func EncodePicks(t *testing.T, options []string) string {
	picks := contract.VoterPicks{VotableIds: options}
	picks_json, err := json.Marshal(picks)
	require.NoError(t, err)

	return string(picks_json)
}
func EncryptOptions(t *testing.T, publicKey *rsa.PublicKey, options []string) string {
	picks_bytes, err := contract.EncryptOAEP(publicKey, []byte(EncodePicks(t, options)))
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(picks_bytes)
}
func SignDataBytes(t *testing.T, privateKey *rsa.PrivateKey, data string) []byte {

	msgHash := sha512.New()
	_, err := msgHash.Write([]byte(data))
	require.NoError(t, err)

	msgHashSum := msgHash.Sum(nil)

	signature_byes, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA512, msgHashSum, nil)
	require.NoError(t, err)

	return signature_byes
}

func SignData(t *testing.T, privateKey *rsa.PrivateKey, data string) string {

	msgHash := sha512.New()
	_, err := msgHash.Write([]byte(data))
	require.NoError(t, err)

	msgHashSum := msgHash.Sum(nil)

	signature_byes, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA512, msgHashSum, nil)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(signature_byes)
}

func CastVote(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]KeyPair, voterId string, options []string) error {

	picks := EncryptOptions(t, keys[voterId].publicKey, options)
	signature := SignData(t, keys[voterId].privateKey, voterId+picks)
	return voterContract.CastVote(transactionContext, voterId, picks, signature)
}

func CastVoteTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]KeyPair, voterId string, optionId string) {
	logSubTest("CastVoteTest")

	err := CastVote(t, transactionContext, voterContract, keys, voterId, []string{optionId})
	require.NoError(t, err)
}

func CastVoteTestWrongUser(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]KeyPair, voterId string) {
	logSubTest("CastVoteTestWrongUser")

	options := []string{"InvalidOption"}
	picks := EncryptOptions(t, keys[voterId].publicKey, options)
	signature := SignData(t, keys[voterId].privateKey, voterId+picks)
	err := voterContract.CastVote(transactionContext, "InvalidUser", picks, signature)
	require.ErrorContains(t, err, contract.ErrNoStateExists.Error())
}
func CastVoteTestWrongOption(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]KeyPair, voterId string) {
	logSubTest("CastVoteTestWrongOption")

	err := CastVote(t, transactionContext, voterContract, keys, voterId, []string{"InvalidOption"})
	require.ErrorContains(t, err, contract.ErrNoStateExists.Error())
}

func CastVoteTestMultiChoice(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]KeyPair, voterId string, optionId1 string, optionId2 string) {
	logSubTest("CastVoteTestMultiChoice")

	err := CastVote(t, transactionContext, voterContract, keys, voterId, []string{optionId1, optionId2})
	if voterContract.SingleChoice {
		require.ErrorContains(t, err, contract.ErrNoVoteIsMoreThanOne.Error()) //Expect error
	} else {
		require.NoError(t, err)
	}
}

func AbstainVoteTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]KeyPair, voterId string) {
	logSubTest("AbstainVoteTest")

	err := CastVote(t, transactionContext, voterContract, keys, voterId, []string{})
	if voterContract.Abstainable {
		require.NoError(t, err)
	} else {
		require.ErrorContains(t, err, contract.ErrNoVoteIsZero.Error()) //Expect error
	}

}

func CastVoteTwiceTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]KeyPair, voterId string, optionId string) {
	logSubTest("CastVoteTwiceTest")

	err := CastVote(t, transactionContext, voterContract, keys, voterId, []string{optionId})
	require.NoError(t, err)

	err = CastVote(t, transactionContext, voterContract, keys, voterId, []string{optionId})
	require.ErrorContains(t, err, contract.ErrAlreadyVoted.Error()) //Revote not allowed

}

func CastVoteWrongKeyTest(t *testing.T, transactionContext *mocks.TransactionContext, voterContract contract.SmartContract, keys map[string]KeyPair, voterId string, wrong_id string, optionId string) {
	logSubTest("CastVoteWrongKeyTest")

	picks := EncryptOptions(t, keys[wrong_id].publicKey, []string{optionId})
	signature := SignData(t, keys[voterId].privateKey, voterId+picks)
	err := voterContract.CastVote(transactionContext, voterId, picks, signature)
	require.ErrorContains(t, err, contract.ErrDecryption.Error()) //wrong pvt key used

	picks = EncryptOptions(t, keys[voterId].publicKey, []string{optionId})
	signature = SignData(t, keys[wrong_id].privateKey, voterId+picks)
	err = voterContract.CastVote(transactionContext, voterId, picks, signature)
	require.ErrorContains(t, err, contract.ErrSignatureValidation.Error()) //wrong pub key used

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
	logTest("TestVote_Anonymous_SingleChoice_Abstainable", 7)

	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3", "User4"}, true, true, true)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTwiceTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	CastVoteWrongKeyTest(t, transactionContext, voterContract, keys, "User2", "User1", "Option2") //Will not refgister

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User2")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User3", "Option1", "Option2") //Will not register

	CastVoteTestWrongUser(t, transactionContext, voterContract, keys, "User1") //Will not register : User1 is for encrypting the options

	CastVoteTestWrongOption(t, transactionContext, voterContract, keys, "User4") //Will not register

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
	logTest("TestVote_Anonymous_SingleChoice_NonAbstainable", 4)

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
	logTest("TestVote_Anonymous_MultiChoice_Abstainable", 4)

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
	logTest("TestVote_Anonymous_MultiChoice_NonAbstainable", 4)

	transactionContext, voterContract, keys, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, true, false, false)
	require.NoError(t, err)

	StateErrorsTests(t, transactionContext, voterContract)

	CastVoteTest(t, transactionContext, voterContract, keys, "User1", "Option1")

	CastVoteTestMultiChoice(t, transactionContext, voterContract, keys, "User2", "Option2", "Option3")

	AbstainVoteTest(t, transactionContext, voterContract, keys, "User3") //will not register

	ExpectedVotingResultTest(t, transactionContext, voterContract, map[string]int{
		"Option1": 1,
		"Option2": 1,
		"Option3": 1,
	})
	ExpectedBallotDetailsTest(t, transactionContext, voterContract, []string{
		"User1",
		"User2",
	})

}

// Test Public, Single-Choice and Abstainable Election
func TestVote_Public_SingleChoice_Abstainable(t *testing.T) {
	logTest("TestVote_Public_SingleChoice_Abstainable", 4)

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
	logTest("TestVote_Public_SingleChoice_NonAbstainable", 4)

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
	logTest("TestVote_Public_MultiChoice_Abstainable", 4)

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
	logTest("TestVote_Public_MultiChoice_NonAbstainable", 4)

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

func EncodePublicKey(publicKey rsa.PublicKey) string {
	publicKey_bytes := x509.MarshalPKCS1PublicKey(&publicKey)
	return base64.StdEncoding.EncodeToString(publicKey_bytes)
}
func DecodePublicKey(t *testing.T, publicKey_base64 string) *rsa.PublicKey {
	publicKey_bytes, err := base64.StdEncoding.DecodeString(publicKey_base64)
	require.NoError(t, err) //Expect No error

	publicKey, err := x509.ParsePKCS1PublicKey(publicKey_bytes)
	require.NoError(t, err) //Expect No error

	return publicKey
}

func AuthUser(t *testing.T, voterContract contract.SmartContract, transactionContext *mocks.TransactionContext, userId string) KeyPair {

	//First authenticate user

	//Generate a key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err) //Expect No error

	publicKey := privateKey.PublicKey
	publicKey_base64 := EncodePublicKey(publicKey)

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

	return KeyPair{publicKey: DecodePublicKey(t, pub), privateKey: privateKey}

}
