package chaincode

import (
    "fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
	"encoding/json"
	"time"
	"encoding/base64"
	"regexp"
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

// user= msps (client id passed in context)

//---------------------------FUNCTIONS---------------------------

// func (s *SmartContract) RegisterBusiness()
 
func(s *SmartContract) VoucherCreated(ctx contractapi.TransactionContextInterface, VoucherID string, SupplierID string, VoucherType string, Hashcode string, TotalValue string, Currency string) error{

	// retrieving id of asset owner (creator)
	OwnerID, err := getClientIdentity(ctx)
	if err != nil {
		return err
	}
	currentTime:= time.Now().UnixMilli() // time.Now.UnixMilli()---- get time in milli seconds

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
		State: "Created",
		
	}
	assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
	state_err := ctx.GetStub().PutState(VoucherID, assetJSON) // new state added

	fmt.Printf("Asset creation returned : %s\n", state_err)

	return state_err


}

func(s *SmartContract) VoucherCancelled(ctx contractapi.TransactionContextInterface, VoucherID string) error{
	cancellingUserID, err := getClientIdentity(ctx)
	if err != nil {
		return err
	}
	VoucherAssetRead, err := s.ReadVoucher(ctx, VoucherID) // voucher asset is read
    if err != nil {
    	return err
    }
	if cancellingUserID != VoucherAssetRead.OwnerID {
		println("You can cancel the voucher only if you are the owner.")
		return fmt.Errorf("You can cancel the voucher only if you are the owner.")
	}

	// ensuring that State of voucher is either "Created" or "Modified"
		State:= VoucherAssetRead.State
		if State == "Created" || State == "Modified"{
					// if the user cancelling the voucher, is the owner ---------------->
					OwnerID:= VoucherAssetRead.OwnerID
					SupplierID:= VoucherAssetRead.SupplierID
					CreatedTime:= VoucherAssetRead.CreatedTime
					UpdatedTime:= time.Now().UnixMilli()
					VoucherType:= VoucherAssetRead.VoucherType
					Hashcode:= VoucherAssetRead.Hashcode
					TotalValue:= VoucherAssetRead.TotalValue
					Currency:= VoucherAssetRead.Currency

					// overwriting original asset with new values
					asset := VoucherAsset{ //creation of asset

						OwnerID: OwnerID,
						SupplierID: SupplierID,
						CreatedTime: CreatedTime,
						UpdatedTime: UpdatedTime,
						VoucherID: VoucherID,
						VoucherType: VoucherType,
						Hashcode: Hashcode,
						TotalValue: TotalValue,
						Currency: Currency,
						State: "Cancelled",
						
					}
					assetJSON, err := json.Marshal(asset)
					if err != nil {
						return err
					}
					state_err := ctx.GetStub().PutState(VoucherID, assetJSON) // asset overridden
					fmt.Printf("Cancellation returned : %s\n", state_err)
					return state_err
	}

	println("You can't cancel when the state is %s", State)
	return fmt.Errorf("You can't cancel when the state is %s", State)

}

func(s *SmartContract) VoucherApproved(ctx contractapi.TransactionContextInterface, VoucherID string) error{

	approvingUserIDFullStr, err := getClientIdentity(ctx)
	if err != nil {
		return err
	}

	approvingUserID:= GetUserid(approvingUserIDFullStr, "x509::CN=",",OU=")



	// EXTRACT "CN" FROM THE APPROVINGUSERID-- in a separate function (try to create a struct-- if the other values are needed in the future)
	VoucherAssetRead, err := s.ReadVoucher(ctx, VoucherID) // voucher asset is read
    if err != nil {
    	return err
    }
	if approvingUserID != VoucherAssetRead.SupplierID {
		println("You can approve the voucher only if you are the supplier.")
		return fmt.Errorf("You can approve the voucher only if you are the supplier.")
	}

	// ensuring that State of voucher is either "Created" or "Modified" in order to approve it
		State:= VoucherAssetRead.State
		if State == "Created" || State == "Modified"{
					// if the user approve the voucher, is the owner ---------------->
					OwnerID:= VoucherAssetRead.OwnerID
					SupplierID:= VoucherAssetRead.SupplierID
					CreatedTime:= VoucherAssetRead.CreatedTime
					UpdatedTime:= time.Now().UnixMilli()
					VoucherType:= VoucherAssetRead.VoucherType
					Hashcode:= VoucherAssetRead.Hashcode
					TotalValue:= VoucherAssetRead.TotalValue
					Currency:= VoucherAssetRead.Currency

					// overwriting original asset with new values
					asset := VoucherAsset{ //creation of asset

						OwnerID: OwnerID,
						SupplierID: SupplierID,
						CreatedTime: CreatedTime,
						UpdatedTime: UpdatedTime,
						VoucherID: VoucherID,
						VoucherType: VoucherType,
						Hashcode: Hashcode,
						TotalValue: TotalValue,
						Currency: Currency,
						State: "Approved",
						
					}
					assetJSON, err := json.Marshal(asset)
					if err != nil {
						return err
					}
					state_err := ctx.GetStub().PutState(VoucherID, assetJSON) // asset overridden
					fmt.Printf("Asset approval returned : %s\n", state_err)
					return state_err
	}

	println("You can't approve when the state is %s", State)
	return fmt.Errorf("You can't approve when the state is %s", State)

}

func(s *SmartContract) VoucherRejected(ctx contractapi.TransactionContextInterface, VoucherID string) error{

	rejectingUserIDFullStr, err := getClientIdentity(ctx)
	if err != nil {
		return err
	}

	rejectingUserID:= GetUserid(rejectingUserIDFullStr, "x509::CN=",",OU=")

	VoucherAssetRead, err := s.ReadVoucher(ctx, VoucherID) // voucher asset is read
    if err != nil {
    	return err
    }
	if rejectingUserID != VoucherAssetRead.SupplierID {
		println("You can reject the voucher only if you are the supplier.")
		return fmt.Errorf("You can reject the voucher only if you are the supplier.")
	}

		// ensuring that State of voucher is either "Created" or "Modified" in order to reject it
		State:= VoucherAssetRead.State
		if State == "Created" || State == "Modified"{
					OwnerID:= VoucherAssetRead.OwnerID
					SupplierID:= VoucherAssetRead.SupplierID
					CreatedTime:= VoucherAssetRead.CreatedTime
					UpdatedTime:= time.Now().UnixMilli()
					VoucherType:= VoucherAssetRead.VoucherType
					Hashcode:= VoucherAssetRead.Hashcode
					TotalValue:= VoucherAssetRead.TotalValue
					Currency:= VoucherAssetRead.Currency

					// overwriting original asset with new values
					asset := VoucherAsset{ //creation of asset

						OwnerID: OwnerID,
						SupplierID: SupplierID,
						CreatedTime: CreatedTime,
						UpdatedTime: UpdatedTime,
						VoucherID: VoucherID,
						VoucherType: VoucherType,
						Hashcode: Hashcode,
						TotalValue: TotalValue,
						Currency: Currency,
						State: "Rejected",
						
					}
					assetJSON, err := json.Marshal(asset)
					if err != nil {
						return err
					}
					state_err := ctx.GetStub().PutState(VoucherID, assetJSON) // asset overridden
					fmt.Printf("Asset rejection returned : %s\n", state_err)
					return state_err
	}

	println("You can't reject when the state is %s", State)
	return fmt.Errorf("You can't reject when the state is %s", State)

}

func(s *SmartContract) VoucherUpdated(ctx contractapi.TransactionContextInterface, VoucherID string, toChange string, newValue string) error{
	// changes in hash or total amount
	// NOTE: Cover both hash and total value updation

	//ensuring that only the owner of the asset can update
	updatingUserID, err := getClientIdentity(ctx)
	if err != nil {
		return err
	}
	VoucherAssetRead, err := s.ReadVoucher(ctx, VoucherID) // voucher asset is read
    if err != nil {
    	return err
    }
	if updatingUserID != VoucherAssetRead.OwnerID {
		println("You can update the voucher only if you are the owner.")
		return fmt.Errorf("You can update the voucher only if you are the owner.")
	}

	State:= VoucherAssetRead.State
	if State == "Created" || State == "Sent Back" || State == "Modified" {
				OwnerID:= VoucherAssetRead.OwnerID
				SupplierID:= VoucherAssetRead.SupplierID
				CreatedTime:= VoucherAssetRead.CreatedTime
				UpdatedTime:= time.Now().UnixMilli()
				VoucherType:= VoucherAssetRead.VoucherType
				Hashcode:= VoucherAssetRead.Hashcode
				TotalValue:= VoucherAssetRead.TotalValue
				Currency:= VoucherAssetRead.Currency

				if toChange=="Hash"{
					asset := VoucherAsset{ 

						OwnerID: OwnerID,
						SupplierID: SupplierID,
						CreatedTime: CreatedTime,
						UpdatedTime: UpdatedTime,
						VoucherID: VoucherID,
						VoucherType: VoucherType,
						Hashcode: newValue,
						TotalValue: TotalValue,
						Currency: Currency,
						State: "Modified",
						
					}
					assetJSON, err := json.Marshal(asset)
					if err != nil {
						return err
					}
					state_err := ctx.GetStub().PutState(VoucherID, assetJSON) // asset overridden
					fmt.Printf("Updation returned : %s\n", state_err)
					return state_err
				} else if toChange=="Value"{

					intermediatenewval, err := strconv.ParseUint(newValue, 10, 32)
					if err !=nil {
						fmt.Println(err)
					}

					newValueInt:= uint(intermediatenewval)
					asset := VoucherAsset{ 

						OwnerID: OwnerID,
						SupplierID: SupplierID,
						CreatedTime: CreatedTime,
						UpdatedTime: UpdatedTime,
						VoucherID: VoucherID,
						VoucherType: VoucherType,
						Hashcode: Hashcode,
						TotalValue: newValueInt,
						Currency: Currency,
						State: "Modified",
						
					}
					assetJSON, err := json.Marshal(asset)
					if err != nil {
						return err
					}
					state_err := ctx.GetStub().PutState(VoucherID, assetJSON) // asset overridden
					fmt.Printf("Updation returned : %s\n", state_err)
					return state_err
				}
			}

		println("You can't update when the state is %s", State)
		return fmt.Errorf("You can't update when the state is %s", State)
}

func(s *SmartContract) VoucherSentBack(ctx contractapi.TransactionContextInterface, VoucherID string) error{

	VoucherAssetRead, err := s.ReadVoucher(ctx, VoucherID) // voucher asset is read
    if err != nil {
    	return err
    }
	requestingUserID, err := getClientIdentity(ctx)
	if err != nil {
		return err
	}
	requestingUser:= GetUserid(requestingUserID, "x509::CN=",",OU=")
	println(requestingUser)
	if requestingUser != VoucherAssetRead.SupplierID{// ensuring that only the supplier can change state to "Sent Back"
		println("You are not a supplier. You cannot perform this action.")
		return fmt.Errorf("You are not a supplier. You cannot perform this action.")
	}
	State:= VoucherAssetRead.State
	if State == "Modified" || State == "Created" { // Verify that the current state of the voucher is created or modified
		OwnerID:= VoucherAssetRead.OwnerID
					SupplierID:= VoucherAssetRead.SupplierID
					CreatedTime:= VoucherAssetRead.CreatedTime
					UpdatedTime:= time.Now().UnixMilli()
					VoucherType:= VoucherAssetRead.VoucherType
					Hashcode:= VoucherAssetRead.Hashcode
					TotalValue:= VoucherAssetRead.TotalValue
					Currency:= VoucherAssetRead.Currency
					asset := VoucherAsset{

						OwnerID: OwnerID,
						SupplierID: SupplierID,
						CreatedTime: CreatedTime,
						UpdatedTime: UpdatedTime,
						VoucherID: VoucherID,
						VoucherType: VoucherType,
						Hashcode: Hashcode,
						TotalValue: TotalValue,
						Currency: Currency,
						State: "Sent Back",
						
					}
					assetJSON, err := json.Marshal(asset)
					if err != nil {
						return err
					}
					state_err := ctx.GetStub().PutState(VoucherID, assetJSON) // asset overridden
					fmt.Printf("Sending back the asset returned : %s\n", state_err)
					return state_err
	}
	println("You can't send back when the state is %s . Requesting User is: %s . Supplier is: %s", State, requestingUser, VoucherAssetRead.SupplierID)
	return fmt.Errorf("You can't send back when the state is %s . Requesting User is: %s . Supplier is: %s", State, requestingUser, VoucherAssetRead.SupplierID)
}

// func(s *SmartContract) UnregisterBusiness()

func(s *SmartContract) ReadVoucher(ctx contractapi.TransactionContextInterface, VoucherID string) (*VoucherAsset, error){

	VoucherAssetJSON, err := ctx.GetStub().GetState(VoucherID)
    if err != nil {
    	return nil, fmt.Errorf("Failed to read from world state: %v", err)
    }
    if VoucherAssetJSON == nil {
    	return nil, fmt.Errorf("A Voucher with ID %s not found.", VoucherID)
    }

    var ReadVoucherAsset VoucherAsset
    err = json.Unmarshal(VoucherAssetJSON, &ReadVoucherAsset)
    if err != nil {
    	return nil, err
	}

	return &ReadVoucherAsset, nil

}


// LOOKUP FUNCTION- GETVOUCHERFORSUPPPLIER() --- gives list of vouchers with a particular supplier

// FUNCTION- GETMYVOUCHERS()- LIST OF VOUCHERS THAT A PARTICULAR OWNER OWNS


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

func GetUserid(Fullstr string, start string, end string) string{

	re := regexp.MustCompile(`CN=([^,]+)`)
    
    match := re.FindStringSubmatch(Fullstr)

	return match[1]

}