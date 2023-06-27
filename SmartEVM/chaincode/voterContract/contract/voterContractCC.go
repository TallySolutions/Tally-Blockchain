package contract

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const Abstained = "_Abstained_"

const VOTE_OPTION_REGISTER_PREFIX = "[VOTE_OPTION]:"
const VOTER_REGISTER_PREFIX = "[VOTE_REGISTER]:"

type SmartContract struct {
	contractapi.Contract
	IsAnonymous  bool
	Initialized  bool
	Abstainable  bool
	SingleChoice bool
}

// NOTE: Write the asset properties in CAMEL CASE- otherwise, chaincode will not get deployed

// Ballot - for each election, each voter will be allotted a single Ballot - one time use only
// Example - for society vote, Ballot is for each flat - so voter it can be the flat-no
//
//	for a group meeting voting, individual assigned one ballot each, here individual id becomes voter id
type Ballot struct {
	//VoterId : for which this ballot is assigned
	VoterId string `json:"Id"`

	//Siignature :  This is base64 encoded digital signature of the VoterId, encrypted using voter specific private key
	Signature string `json:"Signature"`

	//Pub-Key for verify digital signature
	PubKey string `json:"Auth"`

	//CastedVote : Whether vote is casted or not
	Casted bool `json:"Casted"`

	//When the vote is casted
	Timestamp int64 `json:"Timestamp"`

	//Whch options ate voted for - it is kept empty for anonymous voting
	Picks []string `jsno:"Picks"`
}

type VotableOption struct {
	VotableId string `json:"VotableId"`
	Count     int    `json:"Ballots"`
}

// Utility function
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Init Ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, isAnonymous bool, singleChoice bool, abstainable bool) error {
	if s.Initialized {
		return fmt.Errorf("Chaincode already initialized!")
	}
	s.Initialized = true
	s.IsAnonymous = isAnonymous
	s.SingleChoice = singleChoice
	s.Abstainable = abstainable
	if s.Abstainable {
		err := s.AddVotableOption(ctx, Abstained)
		if err != nil {
			s.Initialized = false
			return err
		}
	}
	return nil
}

// Is Votable Option Exists?
func (s *SmartContract) isVotableOptionExists(ctx contractapi.TransactionContextInterface, votableId string) (bool, error) {
	if s.Initialized != true {
		return false, fmt.Errorf("Ledger not inialized!")
	}
	optionBytes, err := ctx.GetStub().GetState(VOTE_OPTION_REGISTER_PREFIX + votableId)
	if err != nil {
		return false, err
	}

	return optionBytes != nil, nil
}

// Is Votable Option Exists?
func (s *SmartContract) isVoterExists(ctx contractapi.TransactionContextInterface, voterId string) (bool, error) {
	if s.Initialized != true {
		return false, fmt.Errorf("Ledger not inialized!")
	}
	optionBytes, err := ctx.GetStub().GetState(VOTER_REGISTER_PREFIX + voterId)
	if err != nil {
		return false, err
	}

	return optionBytes != nil, nil
}

// function to add Votable Option
func (s *SmartContract) AddVotableOption(ctx contractapi.TransactionContextInterface, votableId string) error {

	if s.Initialized != true {
		return fmt.Errorf("Ledger not inialized!")
	}
	fmt.Printf("Adding new votable option: %s\n", votableId)
	//checking if option already added
	optionExists, err := s.isVotableOptionExists(ctx, votableId)
	if err != nil {
		return err
	}
	if optionExists {
		return fmt.Errorf("This votable option is already exist!")
	}

	// if the option does not exist
	votableOption := VotableOption{
		VotableId: votableId,
		Count:     0,
	}
	votableOptionJSON, err := json.Marshal(votableOption)
	if err != nil {
		return err
	}

	fmt.Printf("Creating new asset for this votable id in voting ledger: %s\n", votableId)
	putStateErr := ctx.GetStub().PutState(VOTE_OPTION_REGISTER_PREFIX+votableId, votableOptionJSON) // new state added to the voting ledger
	return putStateErr

}

// function to add Voter
func (s *SmartContract) AddVoter(ctx contractapi.TransactionContextInterface, voterId string) error {

	if s.Initialized != true {
		return fmt.Errorf("Ledger not inialized!")
	}
	fmt.Printf("Adding new voter: %s\n", voterId)
	//checking if voter already added
	optionExists, err := s.isVoterExists(ctx, voterId)
	if err != nil {
		return err
	}
	if optionExists {
		return fmt.Errorf("This voter is already exist!")
	}

	// if the voter does not exist
	ballot := Ballot{
		VoterId:   voterId,
		Signature: "",
		Casted:    false,
		PubKey:    "",
		Timestamp: 0,
	}
	voterJSON, err := json.Marshal(ballot)
	if err != nil {
		return err
	}

	fmt.Printf("Creating new asset for this voter %s\n", voterId)
	putStateErr := ctx.GetStub().PutState(VOTER_REGISTER_PREFIX+voterId, voterJSON) // new state added to the voter ledger
	return putStateErr

}

// function to authorize Voter, this returns a public key specific to this voter's encryption private key
func (s *SmartContract) AuthVoter(ctx contractapi.TransactionContextInterface, voterId string, publicKey_base64 string, signature_base64 string) (string, error) {

	if s.Initialized != true {
		return "", fmt.Errorf("Ledger not inialized!")
	}
	fmt.Printf("Authorizing voter: %s\n", voterId)

	//checking if ballot exists
	ballot, err := s.GetBallot(ctx, voterId)
	if err != nil {
		return "", fmt.Errorf("Unable to get ballot for voter : %v", err)
	}
	if ballot == nil {
		return "", fmt.Errorf("This voter does not exist!")
	}

	// verify signature
	fmt.Println("Verifying signature ...")

	//1. Get bytes from base64 public key
	publicKey_bytes, decodeErr := base64.StdEncoding.DecodeString(publicKey_base64)
	if decodeErr != nil {
		return "", fmt.Errorf("Error in decoding base64 encoded public key: %v", decodeErr)
	}

	//2. Now convert public key bytes into RSA public key
	publicKey, err := x509.ParsePKCS1PublicKey(publicKey_bytes)
	if err != nil {
		return "", fmt.Errorf("Error in parsing public key: %v", err)
	}

	//3. Decode signature
	signature_bytes, SigdecodeErr := base64.StdEncoding.DecodeString(signature_base64)
	if SigdecodeErr != nil {
		return "", fmt.Errorf("Error in decoding base64 encoded signature: %v", SigdecodeErr)
	}

	//4. Verify signatur using public key
	msgHash := sha512.New()
	_, err = msgHash.Write([]byte(voterId))
	if err != nil {
		return "", fmt.Errorf("Error in hashing voter id: %v", err)
	}

	msgHashSum := msgHash.Sum(nil)

	errVerify := rsa.VerifyPSS(publicKey, crypto.SHA512, msgHashSum, signature_bytes, nil)
	if errVerify != nil {
		return "", fmt.Errorf("Error in verifying signature: %v", errVerify)
	}

	// provision key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", fmt.Errorf("Error in generating key: %v", err)
	}

	//1. Get bytes of privatekey
	privateKey_bytes := x509.MarshalPKCS1PrivateKey(privateKey)

	ballot.PubKey = publicKey_base64

	//Encrypt signature

	enc_publicKey := privateKey.PublicKey
	encryptedBytes, err := EncryptOAEP(&enc_publicKey, []byte(signature_base64))
	if err != nil {
		return "", fmt.Errorf("Error in encrypting signature: %v", err)
	}

	encryptedBytes_base64 := base64.StdEncoding.EncodeToString(encryptedBytes)

	//Store base64 encoded of encrypted signature
	ballot.Signature = string(encryptedBytes_base64)

	fmt.Printf("Creating new asset for this voter %s\n", voterId)
	ballotJSON, err := json.Marshal(ballot)
	if err != nil {
		return "", fmt.Errorf("Error in marshalling ballot to json: %v", err)
	}

	putStateErr := ctx.GetStub().PutState(VOTER_REGISTER_PREFIX+voterId, ballotJSON) // updated state added to the voter ledger

	if putStateErr != nil {
		return "", fmt.Errorf("Error in saving ballot: %v", putStateErr)
	}

	//Now store the private key
	privateErr := ctx.GetStub().PutPrivateData(VOTER_REGISTER_PREFIX+voterId, "privateKey", privateKey_bytes)

	if privateErr != nil {
		return "", fmt.Errorf("Error in saving private key: %v", privateErr)
	}

	//1. Get bytes of public key
	enc_publicKey_bytes := x509.MarshalPKCS1PublicKey(&enc_publicKey)

	//2. Base64 encode private key
	enc_publicKey_base64 := base64.StdEncoding.EncodeToString(enc_publicKey_bytes)

	return string(enc_publicKey_base64), nil

}

// function to read vote option
func (s *SmartContract) ReadOption(ctx contractapi.TransactionContextInterface, votableId string) (*VotableOption, error) {
	if s.Initialized != true {
		return nil, fmt.Errorf("Ledger not inialized!")
	}
	votableOptionJSON, err := ctx.GetStub().GetState(VOTE_OPTION_REGISTER_PREFIX + votableId)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if votableOptionJSON == nil {
		return nil, fmt.Errorf("The votable id %s is not present.", votableId)
	}

	var votableOption VotableOption
	err = json.Unmarshal(votableOptionJSON, &votableOption)
	if err != nil {
		return nil, err
	}

	return &votableOption, nil
}

// function to read voter
func (s *SmartContract) GetBallot(ctx contractapi.TransactionContextInterface, voterId string) (*Ballot, error) {
	if s.Initialized != true {
		return nil, fmt.Errorf("Ledger not inialized!")
	}
	voterJSON, err := ctx.GetStub().GetState(VOTER_REGISTER_PREFIX + voterId)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if voterJSON == nil {
		return nil, fmt.Errorf("The voter '%s' is not present.", voterId)
	}

	var ballot Ballot
	err = json.Unmarshal(voterJSON, &ballot)
	if err != nil {
		return nil, err
	}

	return &ballot, nil
}

// function to cast vote, the public key of the ballot's private key must be passed.
func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, voterId string, publicKey_base64 string, votableIds []string) error {
	if s.Initialized != true {
		return fmt.Errorf("Ledger not inialized!")
	}

	if len(votableIds) == 0 {
		if s.Abstainable {
			return s.CastVote(ctx, voterId, publicKey_base64, []string{Abstained})
		}
		return fmt.Errorf("Number of votes to be casted can not be zero.")
	}
	if len(votableIds) > 1 && s.SingleChoice {
		return fmt.Errorf("Number of votes can not be more than one, in case of single choice voting.")
	}

	//checking if ballot exists
	ballot, err := s.GetBallot(ctx, voterId)
	if err != nil {
		return fmt.Errorf("Unable to retrive ballot : %v", err)
	}
	if ballot == nil {
		return fmt.Errorf("This voter does not exist!")
	}

	//Check if user is authorized

	if ballot.PubKey == "" {
		return fmt.Errorf("This voter is not authorized, please authenticate the voter first!")
	}

	if ballot.Casted {
		return fmt.Errorf("This voter is already voted!")
	}

	//Get the ballot's private key
	privateKey_bytes, privateKeyGetErr := ctx.GetStub().GetPrivateData(VOTER_REGISTER_PREFIX+voterId, "privateKey")
	if privateKeyGetErr != nil {
		return fmt.Errorf("Unable to retrive private data : %v", privateKeyGetErr)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKey_bytes)
	if err != nil {
		return fmt.Errorf("Unable to parse private key : %v", err)
	}

	//Decrypt the signature
	signature_bytes, err := base64.StdEncoding.DecodeString(ballot.Signature)
	if err != nil {
		return fmt.Errorf("Unable to decode signature : %v", err)
	}

	decryptedBytes, err := DecryptOAEP(privateKey, signature_bytes)
	if err != nil {
		return fmt.Errorf("Unable to decrypt signature : %v", err)
	}
	signature_base64 := string(decryptedBytes)
	signature_bytes, err = base64.StdEncoding.DecodeString(signature_base64)
	if err != nil {
		return fmt.Errorf("Unable to decode signature : %v", err)
	}

	//Verify signatur using public key
	publicKey_bytes, err := base64.StdEncoding.DecodeString(ballot.PubKey)
	if err != nil {
		return fmt.Errorf("Unable to decode public key : %v", err)
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicKey_bytes)
	if err != nil {
		return fmt.Errorf("Unable to parse public key : %v", err)
	}

	msgHash := sha512.New()
	_, err = msgHash.Write([]byte(voterId))
	if err != nil {
		return fmt.Errorf("Error in hashing voter id: %v", err)
	}

	msgHashSum := msgHash.Sum(nil)

	errVerify := rsa.VerifyPSS(publicKey, crypto.SHA512, msgHashSum, signature_bytes, nil)
	if errVerify != nil {
		return fmt.Errorf("Unable to verify signature : %v", errVerify)
	}

	//Update the timestamp
	ballot.Timestamp = time.Now().UnixNano()
	if s.IsAnonymous {
		ballot.Picks = nil
	} else {
		ballot.Picks = votableIds
	}
	ballot.Casted = true

	for _, optionId := range votableIds {
		votableOption, err := s.ReadOption(ctx, optionId)
		if err != nil {
			return fmt.Errorf("Unable to read option : %v", err)
		}
		votableOption.Count = votableOption.Count + 1 //Vote

		votableOptionJSON, err := json.Marshal(votableOption)
		if err != nil {
			return fmt.Errorf("Unable to marshal option : %v", err)
		}
		updatestate_err := ctx.GetStub().PutState(VOTE_OPTION_REGISTER_PREFIX+optionId, votableOptionJSON)
		if updatestate_err != nil {
			return fmt.Errorf("Unable to update option : %v", updatestate_err)
		}
	}

	//Updare Ballot
	ballotJSON, err := json.Marshal(ballot)
	if err != nil {
		return fmt.Errorf("Unable to marshal ballot : %v", err)
	}
	updatestate_err := ctx.GetStub().PutState(VOTER_REGISTER_PREFIX+voterId, ballotJSON)
	if updatestate_err != nil {
		return fmt.Errorf("Unable to update ballot : %v", updatestate_err)
	}

	return nil
}

// GetAllOptions returns all voting options found in world state
func (s *SmartContract) GetAllOptions(ctx contractapi.TransactionContextInterface) ([]*VotableOption, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}

	if resultsIterator == nil {
		return nil, fmt.Errorf("No result iterator found!")
	}

	defer resultsIterator.Close()

	var assets []*VotableOption
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(queryResponse.Key, VOTE_OPTION_REGISTER_PREFIX) {
			var asset VotableOption
			err = json.Unmarshal(queryResponse.Value, &asset)
			if err != nil {
				return nil, err
			}
			assets = append(assets, &asset)
		}
	}

	return assets, nil
}

// GetAllBallots returns all voting ballots found in world state
func (s *SmartContract) GetAllBallots(ctx contractapi.TransactionContextInterface) ([]*Ballot, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}

	if resultsIterator == nil {
		return nil, fmt.Errorf("No result iterator found!")
	}

	defer resultsIterator.Close()

	var assets []*Ballot
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(queryResponse.Key, VOTER_REGISTER_PREFIX) {
			var asset Ballot
			err = json.Unmarshal(queryResponse.Value, &asset)
			if err != nil {
				return nil, err
			}
			assets = append(assets, &asset)
		}
	}

	return assets, nil
}

func EncryptOAEP(public *rsa.PublicKey, msg []byte) ([]byte, error) {
	hash := sha512.New()
	random := rand.Reader

	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], nil)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func DecryptOAEP(private *rsa.PrivateKey, msg []byte) ([]byte, error) {
	hash := sha512.New()
	random := rand.Reader

	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], nil)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
