package chaincode

import (
    "fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
	"encoding/json"
	"time"
	"encoding/base64"
)


type SmartContract struct {
    contractapi.Contract
}

type VoucherAsset struct {
	OwnerID string `json:"OwnerID"`
	SupplierID string `json:"SupplierID"`
	CreatedTime int64 `json:"CreatedTime"`
	UpdatedTime int64 `json:"UpdatedTime"`
	VoucherID string `json:"VoucherID"`
	VoucherType string `json:"VoucherType"`
	Hashcode string `json:"Hashcode"`
	TotalValue uint `json"TotalValue"`
	Currency string `json:"Currency"`
	State string `json:"State"`
} 

// should hashcode be [64]byte?
// time.Now()-- returns current time as type time.Time

// user= msps (client id passed in context)

// id= voucher id

// time.Now.UnixMilli()---- get time in milli seconds

//---------------------------FUNCTIONS---------------------------


// supplier id-- passed as a param for me (only this or reject or approve or send back)
// only the owner id -- can change the state

// func (s *SmartContract) RegisterBusiness()
 
func(s *SmartContract) VoucherCreated(ctx contractapi.TransactionContextInterface, VoucherID string, SupplierID string, VoucherType string, Hashcode string, TotalValue string, Currency string, State string ) error{

	// retrieving id of asset owner (creator)
	OwnerID, err := getClientIdentity(ctx)
	if err != nil {
		return err
	}
	currentTime:= time.Now().UnixMilli()

	Value, err := strconv.ParseUint(TotalValue, 10, 32)
    if err !=nil {
    	fmt.Println(err)
    }
    TotalValueStr := uint(Value)

	asset := VoucherAsset{ //creation of asset

		OwnerID: OwnerID,
		SupplierID: SupplierID,
		CreatedTime: currentTime,
		UpdatedTime: currentTime,
		VoucherID: VoucherID,
		VoucherType: VoucherType,
		Hashcode: Hashcode,
		TotalValue: TotalValueStr,
		Currency: Currency,
		State: State,
		
	}
	assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		state_err := ctx.GetStub().PutState(VoucherID, assetJSON) // new state added

		fmt.Printf("Asset creation returned : %s\n", state_err)

		return state_err


}

func(s *SmartContract) VoucherCancelled(ctx contractapi.TransactionContextInterface){

}

func(s *SmartContract) VoucherApproved(ctx contractapi.TransactionContextInterface){

}

func(s *SmartContract) VoucherRejected(ctx contractapi.TransactionContextInterface){

}

func(s *SmartContract) VoucherUpdated(ctx contractapi.TransactionContextInterface){

}

func(s *SmartContract) VoucherSentBack(ctx contractapi.TransactionContextInterface){

}

// func(s *SmartContract) UnregisterBusiness()


func getClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil     // returns clientID as a string

}