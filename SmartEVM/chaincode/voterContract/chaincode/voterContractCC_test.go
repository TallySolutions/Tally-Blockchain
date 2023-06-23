package chaincode_test

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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
