package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
	"tallysolutions.com/SmartEVM/chaincode/voterContract/chaincode"
	"tallysolutions.com/SmartEVM/chaincode/voterContract/chaincode/mocks"
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

	voterContract := chaincode.SmartContract{}
	err := voterContract.InitLedger(transactionContext)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = voterContract.InitLedger(transactionContext)
	require.EqualError(t, err, "failed inserting key")
}

func TestAddVotableOption(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	//Running without initialization
	voterContract := chaincode.SmartContract{}
	err := voterContract.AddVotableOption(transactionContext, "Option1")
	require.Error(t, err) //Expect error

	//Initialize
	err = voterContract.InitLedger(transactionContext)
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

func TestCastVote(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	//Running without initialization
	voterContract := chaincode.SmartContract{}
	err := voterContract.CastVote(transactionContext, "User1", []string{"Option1"})
	require.Error(t, err) //Expect error

	//Initialize
	err = voterContract.InitLedger(transactionContext)
	require.NoError(t, err) //Expect No error

	//Running without any options
	err = voterContract.CastVote(transactionContext, "User1", []string{"Option1"})
	require.Error(t, err) //Expect error - as no option Option1 present yet

	//Setup the assets
	expectedAsset := &chaincode.VotableOption{VotableId: "Option1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	//Set the asset for iterator as well
	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.GetStateByRangeReturns(iterator, nil)

	//Run with a set option = Option1
	err = voterContract.CastVote(transactionContext, "User1", []string{"Option1"})
	require.NoError(t, err) //Expect no error

}
