package contract_test

import (
	"fmt"
	"testing"

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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
	require.EqualError(t, err, "Chaincode already initialized!")
}

func TestInitLedgerStateFailure(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	voterContract := contract.SmartContract{}

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err := voterContract.InitLedger(transactionContext, true, true, true)
	require.EqualError(t, err, "failed inserting key")
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
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset."))
	err = voterContract.AddVotableOption(transactionContext, "Option1")
	require.EqualError(t, err, "unable to retrieve asset.")
}

var states map[string][]byte

type PrivateData map[string][]byte

var collections map[string]PrivateData

func GetStateStub(id string) ([]byte, error) {
	if id == "_state_error_" {
		return nil, fmt.Errorf("Unable to retrieve value from state")
	}
	return states[id], nil
}
func AddStateStub(id string, value []byte) error {
	if id == "_state_error_" {
		return fmt.Errorf("Unable to set value from state")
	}
	states[id] = value

	return nil
}

func PutPrivateDataStub(collection string, key string, value []byte) error {
	if collection == "_state_error_" {
		return fmt.Errorf("Unable to set value from state")
	}
	data := PrivateData{}
	data[key] = value
	collections[collection] = data

	return nil
}

func GetPrivateDataStub(collection string, key string) ([]byte, error) {
	if collection == "_state_error_" {
		return nil, fmt.Errorf("Unable to get value from state")
	}

	return collections[collection][key], nil
}

func SetupVote(t *testing.T, options []string, voters []string, anonymous bool, singlechoice bool, abstainable bool) (*mocks.TransactionContext, contract.SmartContract, error) {
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

	//Initialize
	err := voterContract.InitLedger(transactionContext, anonymous, singlechoice, abstainable)
	require.NoError(t, err)

	//Set options
	for _, optionId := range options {
		err = voterContract.AddVotableOption(transactionContext, optionId)
		require.NoError(t, err)
	}

	//Set options
	for _, voterId := range voters {
		err = voterContract.AddVoter(transactionContext, voterId)
		require.NoError(t, err)
	}

	return transactionContext, voterContract, err
}

// Test Anonymous, Single-Choice and Abstainable Election
func TestVote_Anonymous_SingleChoice_Abstainable(t *testing.T) {

	transactionContext, voterContract, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, true, true, true)
	require.NoError(t, err)

	err = voterContract.CastVote(transactionContext, "User1", "", []string{"option1", "option2"})
	require.EqualError(t, err, "Number of votes can not be more than one, in case of single choice voting.") //Expect error : No voting choice

	pubkey := AuthUser(t, voterContract, transactionContext, "User1")

	err = voterContract.CastVote(transactionContext, "User1", pubkey, []string{"Option1"})
	require.NoError(t, err)
}

// Test Anonymous, Single-Choice and Not-Abstainable Election
func TestVote_Anonymous_SingleChoice_NonAbstainable(t *testing.T) {
	transactionContext, voterContract, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, true, true, false)
	require.NoError(t, err)

	err = voterContract.CastVote(transactionContext, "User1", "", []string{"option1", "option2"})
	require.EqualError(t, err, "Number of votes can not be more than one, in case of single choice voting.") //Expect error : multiple choice

	err = voterContract.CastVote(transactionContext, "User1", "", []string{})
	require.EqualError(t, err, "Number of votes to be casted can not be zero.") //Expect error : No voting choice

}

// Test Anonymous, Multi-Choice and Not-Abstainable Election
func TestVote_Anonymous_MultiChoice_NonAbstainable(t *testing.T) {
	transactionContext, voterContract, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, true, false, false)
	require.NoError(t, err)

	err = voterContract.CastVote(transactionContext, "User1", "", []string{})
	require.EqualError(t, err, "Number of votes to be casted can not be zero.") //Expect error : No voting choice
}

// Test Public, Single-Choice and Abstainable Election
func TestVote_Public_SingleChoice_Abstainable(t *testing.T) {
	transactionContext, voterContract, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, false, true, true)
	require.NoError(t, err)

	err = voterContract.CastVote(transactionContext, "User1", "", []string{"option1", "option2"})
	require.EqualError(t, err, "Number of votes can not be more than one, in case of single choice voting.") //Expect error : Multiple choices
}

// Test Public, Single-Choice and Not-Abstainable Election
func TestVote_Public_SingleChoice_NonAbstainable(t *testing.T) {
	transactionContext, voterContract, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, false, true, false)
	require.NoError(t, err)

	err = voterContract.CastVote(transactionContext, "User1", "", []string{})
	require.EqualError(t, err, "Number of votes to be casted can not be zero.") //Expect error : No voting choice

	err = voterContract.CastVote(transactionContext, "User1", "", []string{"option1", "option2"})
	require.EqualError(t, err, "Number of votes can not be more than one, in case of single choice voting.") //Expect error :multiple choices

}

// Test Public, Multi-Choice and Not-Abstainable Election
func TestVote_Public_MultiChoice_NonAbstainable(t *testing.T) {
	transactionContext, voterContract, err := SetupVote(t, []string{"Option1", "Option2", "Option3"}, []string{"User1", "User2", "User3"}, false, false, false)
	require.NoError(t, err)

	err = voterContract.CastVote(transactionContext, "User1", "", []string{})
	require.EqualError(t, err, "Number of votes to be casted can not be zero.") //Expect error : No voting choice

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

	pub, err := voterContract.AuthVoter(transactionContext, userId, publicKey_base64, signature_base64)
	require.NoError(t, err) //Expect No error

	return pub

}
