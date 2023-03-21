package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type OwnerAsset struct {
	OwnerID   string `json:"OwnerID"`
	OwnerName string `json:"OwnerName"`
	IsActive  bool   `json:"IsActive"`
}

const Prefix = "Owner:"

func (s *SmartContract) IsOwnerActive(ctx contractapi.TransactionContextInterface, OwnerID string) (bool, error) {
	// returns boolean for owner status
	ownerJSON, err := ctx.GetStub().GetState(Prefix + OwnerID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	var owner OwnerAsset
	err2 := json.Unmarshal([]byte(ownerJSON), &owner)

	if err2 != nil {
		return false, fmt.Errorf("failed conversion to JSON in checking active status: %v", err2)
	}

	if owner.IsActive {
		fmt.Printf("Owner exists returned : %t\n", owner.IsActive)
		return true, nil
	} else {
		fmt.Printf("Owner exists returned : %t\n", owner.IsActive)
		return false, nil
	}
}

func (s *SmartContract) MakeOwnerActive(ctx contractapi.TransactionContextInterface, OwnerID string) error {
	ownerJSON, err := ctx.GetStub().GetState(Prefix + OwnerID)
	if err != nil {
		return fmt.Errorf("Failed to read from world state: %v", err)
	}
	var owner OwnerAsset
	err2 := json.Unmarshal([]byte(ownerJSON), &owner)
	if err2 != nil {
		return fmt.Errorf("Failed conversion to JSON in checking active status: %v", err2)
	}
	owner.IsActive = true
	updatedOwnerJSON, err := json.Marshal(owner)
	if err != nil {
		return fmt.Errorf("failed to marshal updated owner: %v", err)
	}

	return ctx.GetStub().PutState(Prefix+OwnerID, updatedOwnerJSON)
}

func (s *SmartContract) MakeOwnerInactive(ctx contractapi.TransactionContextInterface, OwnerID string) error {
	ownerJSON, err := ctx.GetStub().GetState(Prefix + OwnerID)
	if err != nil {
		return fmt.Errorf("Failed to read from world state: %v", err)
	}
	var owner OwnerAsset
	err2 := json.Unmarshal([]byte(ownerJSON), &owner)
	if err2 != nil {
		return fmt.Errorf("Failed conversion to JSON in checking active status: %v", err2)
	}
	owner.IsActive = false
	updatedOwnerJSON, err := json.Marshal(owner)
	if err != nil {
		return fmt.Errorf("failed to marshal updated owner: %v", err)
	}

	return ctx.GetStub().PutState(Prefix+OwnerID, updatedOwnerJSON)
}

func (s *SmartContract) OwnerExistence(ctx contractapi.TransactionContextInterface, OwnerID string) (bool, error) {
	ownerJSON, err := ctx.GetStub().GetState(Prefix + OwnerID)

	if err != nil {

		return false, fmt.Errorf("Failed to read from world state: %v", err)

	}

	return ownerJSON != nil, nil
}

func (s *SmartContract) RegisterOwner(ctx contractapi.TransactionContextInterface, OwnerID string, Name string) error {
	// NAME SHOULD BE PASSED AS A PARAMETER
	// id generation happens in this function- on creation of an owner
	ownerexists, err := s.OwnerExistence(ctx, OwnerID) // TODO: FIND OUT HOW TO CHECK FOR OWNER NAME INSTEAD OF ID (gen ID when registering owner- not here)

	if err != nil {
		return err
	}

	// if owner exists already

	if ownerexists {
		// now there are 2 possible scenarios- active, the other is inactive
		owneractive, err := s.IsOwnerActive(ctx, OwnerID)
		if err != nil {
			return err
		}
		if owneractive { // if owner is active i.e. existing and active, a statement is returned that the user is already registered
			fmt.Printf("ERRROR : Owner is already registered!")
			return fmt.Errorf("Owner is already registered")
		} else {
			// if the owner is not active, they are made active
			err := s.MakeOwnerActive(ctx, OwnerID)
			if err != nil {
				return fmt.Errorf("error in changing owner's status")
			}
			fmt.Printf("Owner is active")
			return nil
		}
	}

	// OWNER DOES NOT EXIST- So, now we will create a new owner, and register them too
	owner := OwnerAsset{
		OwnerName: Name,
		// OwnerID:   uuid.New().String(),
		OwnerID:  OwnerID,
		IsActive: true,
	}
	ownerJSON, err := json.Marshal(owner)
	if err != nil {
		return err
	}
	state_err := ctx.GetStub().PutState(Prefix+OwnerID, ownerJSON) // new state added
	fmt.Printf("Owner creation returned : %s\n", state_err)

	return state_err

}

func (s *SmartContract) UnregisterOwner(ctx contractapi.TransactionContextInterface, OwnerID string) error {

	ownerexists, err := s.OwnerExistence(ctx, OwnerID)

	if err != nil {
		return err
	}
	if ownerexists {
		owneractive, err := s.IsOwnerActive(ctx, OwnerID)
		if err != nil {
			return err
		}
		if owneractive { // if owner is active i.e. existing and active, the owner is made inactive
			err := s.MakeOwnerInactive(ctx, OwnerID)
			if err != nil {
				return fmt.Errorf("error in changing owner's status")
			}
			fmt.Printf("Owner is active")
			return nil
		} else {
			fmt.Printf("ERRROR : Owner is already unregistered!")
			return fmt.Errorf("Owner is already unregistered")
		}
	}

	return nil
}

func (s *SmartContract) GetAllOwners(ctx contractapi.TransactionContextInterface) ([]*OwnerAsset, error) {
	iteratorVar, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {

		return nil, err

	}

	defer iteratorVar.Close()

	var owners []*OwnerAsset
	var ownerCount = 0

	for iteratorVar.HasNext() {
		queryResponse, err := iteratorVar.Next()
		if err != nil {
			return nil, err
		}

		var owner OwnerAsset
		err = json.Unmarshal(queryResponse.Value, &owner)
		if err != nil {
			return nil, err
		}
		owners = append(owners, &owner)
		ownerCount++
	}

	if ownerCount > 0 {
		return owners, nil
	} else {
		return nil, fmt.Errorf("No owners found")
	}

}
