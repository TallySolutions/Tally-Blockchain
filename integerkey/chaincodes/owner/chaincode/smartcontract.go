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
	ownerJSON, err := ctx.GetStub().GetState(Prefix + Name)
	if err != nil {

		return false, fmt.Errorf("failed to read from world state: %v", err)

	}

	var owner OwnerAsset
	err := json.Unmarshal([]byte(ownerJSON), &owner)

		if err != nil {
			fmt.Errorf("failed conversion to JSON in checking active status")
			return nil, err
		}

	if owner.isActive{

		fmt.Printf("Owner exists returned : %t\n", owner.isActive)
		return true,nil
	} else {
			fmt.Printf("Owner exists returned : %t\n", owner.isActive)
			return false, nil
	}

}

func (s *SmartContract) MakeOwnerActive(ctx contractapi.TransactionContextInterface , Name string) error{

	ownerJSON, err := ctx.GetStub().GetState( Prefix + Name)
		if err != nil {
			return fmt.Errorf("Failed to read from world state: %v", err)
		}
	var owner OwnerAsset
	err := json.Unmarshal([]byte(ownerJSON), &owner)
		if err != nil {
			fmt.Errorf("Failed conversion to JSON in checking active status")
			return err
		}
	owner.isActive = true
	return nil
}

func (s *SmartContract) MakeOwnerInactive(ctx contractapi.TransactionContextInterface , Name string) error{

	ownerJSON, err := ctx.GetStub().GetState(Prefix + Name)
		if err != nil {
			return fmt.Errorf("Failed to read from world state: %v", err)
		}
	var owner OwnerAsset
	err := json.Unmarshal([]byte(ownerJSON), &owner)
		if err != nil {
			fmt.Errorf("Failed conversion to JSON in checking active status")
			return err
		}
	owner.isActive = false
	return nil
}


func (s *SmartContract) OwnerExistence(ctx contractapi.TransactionContextInterface, Name string) (bool, error) {
	ownerJSON, err := ctx.GetStub().GetState(Prefix + Name)

	if err != nil {

		return false, fmt.Errorf("Failed to read from world state: %v", err)

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
				owneractive, err := s.IsOwnerActive(ctx, Name)
				if err != nil {
					return err
				}
				if owneractive{		// if owner is active i.e. existing and active, a statement is returned that the user is already registered
						fmt.Printf("ERRROR : Owner is already registered!")
						return fmt.Errorf("Owner is already registered")
				} else{	
						// if the owner is not active, they are made active
						err := s.MakeOwnerActive(ctx, Name)
						if err != nil {
							return fmt.Errorf("error in changing owner's status")
						}
						fmt.Printf("Owner is active")
						return nil
				}
		} 
		
		// OWNER DOES NOT EXIST- So, now we will create a new owner, and register them too
		owner := Owner{
			OwnerName: Name,
			OwnerID: uuid.New().String(),
			isActive: true
		}
		ownerJSON, err := json.Marshal(owner)
		if err != nil {
			return err
		}
		state_err := ctx.GetStub().PutState(Prefix + Name, ownerJSON) // new state added
		fmt.Printf("Owner creation returned : %s\n", state_err)

		return state_err


}