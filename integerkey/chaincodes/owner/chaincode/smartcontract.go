package chaincode


// go get github.com/google/uuid -- to be done in terminal



import(
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/google/uuid"
)


type SmartContract struct{
	contractapi.Contract
}

type OwnerAsset struct{
	OwnerID string `json:OwnerID`
	OwnerName string `json:OwnerName`
	isActive bool `json:isActive`
}

const Prefix = "Owner: "




func (s *SmartContract) IsOwnerActive(ctx contractapi.TransactionContextInterface, Name string) (bool,error){
	// returns boolean for owner status
}



func (s *SmartContract) OwnerExistence(ctx contractapi.TransactionContextInterface, Name string) (bool, error) {
	ownerJSON, err := ctx.GetStub().GetState(Prefix + Name)

	if err != nil {

		return false, fmt.Errorf("failed to read from world state: %v", err)

	}

	return assetJSON != nil, nil
}

func (s *SmartContract) RegisterOwner(ctx TransactionContextInterface, Name string) error {
	// id generation happens in this function- on creation of an owner
	ownerexists, err := s.OwnerExistence(ctx, Name)

	if err != nil {
		return err
	}

	// if owner exists already
	if ownerexists {
		// now there are 2 possible scenarios- active, the other is inactive
	}

}