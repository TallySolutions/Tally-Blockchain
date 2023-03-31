package chaincode


import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "strings"
    "strconv"
    // "github.com/google/uuid"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}
type Asset struct {
    Name  string `json:"Name"`
    Value uint   `json:"Value"`
    OwnerID string `json:"OwnerID"`
}



// ADD REQUESTED OWNER to struct- WHEN IMPLMENTING FUNC


type OwnerAsset struct {
	OwnerID   string `json:"OwnerID"`
	OwnerName string `json:"OwnerName"`
	IsActive  bool   `json:"IsActive"`
}



const Prefix = "Key: "
const OwnerPrefix="Owner: "

// function that takes input as context of transaction and the name of the key, returns boolean value that implies whether the asset exists or not, otherwise- an error
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, Name string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(Prefix + Name)
    if err != nil {
    return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
}


func (s *SmartContract) GetAssetValue(ctx contractapi.TransactionContextInterface, Name string)(string, error){
	//returns ownerID to the app that calls it
	assets_list, err := s.GetAllAssets(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	var asset *Asset
	for _, iteratorVar := range assets_list{
		if iteratorVar.Name == Name{
			asset = iteratorVar
			break
		}
	}
	//check for existence 
	// found and retrieved the matching owner, now we have to return the id of the owner
	return string(asset.Value), nil
	
}



// function to create an asset. Input= transaction context, name of the key to be created. Creates new asset if an asset with the name given does not exist
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, Name string) error {


    OwnerID, err := submittingClientIdentity(ctx)
	if err != nil {
		return err
	}
    // get id from func
    exists, err := s.AssetExists(ctx, Prefix + Name) // exists-> boolean value, err-> can be nil or the error, if present

    fmt.Printf("Asset exists returned : %t, %s\n", exists, err)

    if err != nil {
    return err
                }
        if exists {
            return fmt.Errorf("the asset %s already exists", Prefix + Name)
            }

            asset := Asset{ //creation of asset
                Name: Name,
                Value: 0,
                OwnerID: OwnerID,
            }
            assetJSON, err := json.Marshal(asset)
            if err != nil {
            return err
            }

        state_err := ctx.GetStub().PutState(Prefix + Name, assetJSON) // new state added

        fmt.Printf("Asset creation returned : %s\n", state_err)

        return state_err
}

// ReadAsset returns the asset stored in the world state with given Name.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, Name string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(Prefix + Name)
    if err != nil {
    return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
    return nil, fmt.Errorf("the asset %s does not exist", Name)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
    return nil, err
}

return &asset, nil
}



func (s *SmartContract )GetAssetsPagination(ctx contractapi.TransactionContextInterface, startname string, endname string, bookmark string) ([] *Asset, error) {

    // NOTE: BOOKMARK HAS TO BE SENT AS AN EMPTY STRING WHEN SENT AS A PARAMETER
    // pageSizeInt, e := strconv.Atoi(pageSize)
    // if e != nil {
    // 	return nil, e
    //   }
    pageSizeInt := 5
                   iteratorVar, midvar, err:= ctx.GetStub().GetStateByRangeWithPagination(Prefix + startname, Prefix + endname, int32(pageSizeInt), bookmark)
    if err !=nil && midvar!=nil {
    return nil, err
}
defer iteratorVar.Close()


    var assets []*Asset

    for iteratorVar.HasNext() {
        queryResponse, err := iteratorVar.Next()
        if err != nil {
        return nil, err
    }

    var asset Asset
    err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
        return nil, err
    }
    assets = append(assets, &asset)
    }

    return assets, nil

}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([] *Asset, error) {

    iteratorVar, err := ctx.GetStub().GetStateByRange("","")   // TRY RANGE PARAMETERS , other getstateby.... (rows etc.)
    if err !=nil {
    return nil, err
}
defer iteratorVar.Close()

    var assets []*Asset

    var assetCount = 0
    for iteratorVar.HasNext() {
        queryResponse, err := iteratorVar.Next()
        if err != nil {
        return nil, err
    }

    var asset Asset
    err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
        return nil, err
    }
    assets = append(assets, &asset)
                 assetCount++
    }

    if assetCount > 0 {
    return assets, nil
} else {
    return nil, fmt.Errorf("No assets found")
    }

}



// IncreaseAsset increases the value of the asset by the specified value- with certain limits
func (s *SmartContract) IncreaseAsset(ctx contractapi.TransactionContextInterface, Name string, incrementValue string) (*Asset, error) {
    // NOTE: incrementValue is a string because SubmitTransaction accepts string parameters as input parameters
    // accepting owner because we will be OVERWRITING the asset
    asset_read, err := s.ReadAsset(ctx, Name) // asset is read
    if err != nil {
    return nil, err
}


intermediateUpdateval, err := strconv.ParseUint(incrementValue, 10, 32)
    if err !=nil {
    fmt.Println(err)
    }
    incrementValueuInt := uint(intermediateUpdateval)
    newValue := uint(asset_read.Value) + incrementValueuInt

    if newValue > 20 {
    return nil, fmt.Errorf("You cannot have a value more than 20.")
    }

    // overwriting original asset with new value
    asset := Asset {
        Name:  Name,
        Value: newValue,
        // OwnerID: ownerID,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
    return nil, err
}

updatestate_err := ctx.GetStub().PutState(Prefix + Name, assetJSON)
                       fmt.Printf("Increasing asset value returned the following: %s ", updatestate_err)

                       return &asset, nil
}

// DecreaseAsset decreases the value of the asset by the specified value
func (s *SmartContract) DecreaseAsset(ctx contractapi.TransactionContextInterface, Name string, decrementValue string) (*Asset, error) {
    asset_read, err := s.ReadAsset(ctx, Name)
    if err != nil {
    return nil, err
    }

intermediateval, err := strconv.ParseUint(decrementValue, 10, 32)
    if err !=nil {
    fmt.Println(err)
    }
    decrementValueuInt := uint(intermediateval)
    if decrementValueuInt > uint(asset_read.Value) {
        return nil, fmt.Errorf("You cannot decrement value to less than 0.")
    }
    newValue := uint(asset_read.Value) - decrementValueuInt


                // overwriting original asset with new value
    asset := Asset {
        Name:  Name,
        Value: newValue,
        // OwnerID: ownerID,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
    return nil, err
    }

updatestate_Err := ctx.GetStub().PutState(Prefix + Name, assetJSON)
                       fmt.Printf("After decreasing asset value: %s", updatestate_Err)

                       return &asset , nil
}



func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, Name string) (*Asset, error) {
    // asset_read, err := s.ReadAsset(ctx, Name)
    // if err != nil {
    // return nil, err
    // }
    // overwriting original asset with new owner

    // finding the asset with Name provided as param
    // assetJSON, err:= ctx.GetStub().GetState(Prefix + Name)
    // if err != nil {
    //     return nil, err
    // }
    // println(assetJSON)
    // var asset Asset
    // err = json.Unmarshal(assetJSON, &asset)
    // if err != nil {
    // return nil, err
    // }

    newOwnerID, err := submittingClientIdentity(ctx)
	if err != nil {
		return nil,err
	}


    // REMOVE OWNERID AS PARAM-- GET OWNERID FROM CONTEXT(USER CALLING IT)

    iteratorVar, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {

		return nil, err

	}

	defer iteratorVar.Close()

	for iteratorVar.HasNext() {
		queryResponse, err := iteratorVar.Next()
		if err != nil {
			return nil, err
		}
        println(queryResponse.Value)
    }

    asset, err:= s.ReadAsset(ctx, Name)
    if err != nil {
        return nil, err
    }
    println(asset.OwnerID)
    ownerassetJSON, err:= ctx.GetStub().GetState(OwnerPrefix + asset.OwnerID)
    if err != nil {
        return nil, err
    }
    println(string(ownerassetJSON))
    var ownerasset OwnerAsset
    err = json.Unmarshal(ownerassetJSON, &ownerasset)
    if err != nil {
    return nil, err
    }
    currownerName:= ownerasset.OwnerName
    println("Current owner Name:"+ currownerName)

    // we have retrived the current owner name.. now we have to verify if it is active
    if !ownerasset.IsActive{
        return nil, fmt.Errorf("PROBLEM: %s", "not active")
    }

    // if current owner is active, continue

    println("New owner ID:" + newOwnerID)
    
    // overwriting current asset with new owner id
    val_AssetInt:= asset.Value
    new_asset := Asset {
        Name:  Name,
        Value: val_AssetInt,
        OwnerID: newOwnerID,
    }

    assetJSON, err := json.Marshal(new_asset)
    if err != nil {
    return nil, err
    }

    ctx.GetStub().PutState(Prefix + Name, assetJSON)

    return &new_asset , nil
}



// DeleteAsset deletes the state from the ledger
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, name string) error {
    exists, err := s.AssetExists(ctx, name)
    if err != nil {
    return err
}
if !exists {
    return fmt.Errorf("the asset %s does not exist", name)
    }

    delop:= ctx.GetStub().DelState(Prefix + name)
    fmt.Printf("Message received on deletion: %s", delop)
    return nil
}


// trasnfer( asset id, destination owner)

func (s *SmartContract) IsOwnerActive(ctx contractapi.TransactionContextInterface, Name string) (string, error) {
	// returns boolean for owner status
	owners_list, err := s.GetAllOwners(ctx)
	if err != nil {
		return "false", fmt.Errorf("failed to read from world state: %v", err)
	}

	var owner *OwnerAsset
	for _, iteratorVar := range owners_list{
		if iteratorVar.OwnerName == Name{
			owner= iteratorVar
			break
		}
		return "owner does not exist", nil
	}
	if owner.IsActive{
		return "true", nil
	} else{
		return "false", nil
	}
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

		if strings.HasPrefix(queryResponse.Key, OwnerPrefix) {
			var owner OwnerAsset
			err = json.Unmarshal(queryResponse.Value, &owner)
			if err != nil {
				return nil, err
			}
			owners = append(owners, &owner)
			ownerCount++
		}
	}

	if ownerCount > 0 {
		return owners, nil
	} else {
		return []*OwnerAsset{}, nil
	}

}



func submittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
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






// REQUEST TRANSFER- req transanction

// two modes- automatic or requires approval(based on condition- asset value<10 or >=10)








// APPROVE TRANSFER- for the user who owns the asset










// AFTER THESE TWO, READ THROUGH ABAC        (access control)