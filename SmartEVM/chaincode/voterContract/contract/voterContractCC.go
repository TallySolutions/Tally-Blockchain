package contract

import (
	"encoding/json"
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
	PublicKey    string
	PrivateKey   string
	log          CCLog
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

type VoterPicks struct {
	VotableIds []string `json:"VotableIds"`
}

// Init Ledger
// Must supply base64 encoded public keys of the admin users
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, admin_public_key string, isAnonymous bool, singleChoice bool, abstainable bool, Debug bool) (string, error) {
	if s.Initialized {
		return "", Error(ErrCCAlreadyInitialized)
	}
	s.Initialized = true
	s.IsAnonymous = isAnonymous
	s.SingleChoice = singleChoice
	s.Abstainable = abstainable
	s.PublicKey = admin_public_key
	s.log = CCLog{}
	if Debug {
		s.log.PrintDebug = true
	}
	if s.Abstainable {
		err := s.AddVotableOption(ctx, Abstained)
		if err != nil {
			s.Initialized = false
			return "", err
		}
	}

	// provision key pair fo admin
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", WrapError(ErrKeyGeneration, err)
	}

	// Get bytes of privatekey
	privateKey_bytes := x509.MarshalPKCS1PrivateKey(privateKey)
	s.PrivateKey = base64.StdEncoding.EncodeToString(privateKey_bytes)

	ppublicKey_bytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	return base64.StdEncoding.EncodeToString(ppublicKey_bytes), nil
}

// Is Votable Option Exists?
func (s *SmartContract) IsVotableOptionExists(ctx contractapi.TransactionContextInterface, votableId string) (bool, error) {
	if s.Initialized != true {
		return false, Error(ErrLedgerNotInitialized)
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
		return false, Error(ErrLedgerNotInitialized)
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
		return Error(ErrLedgerNotInitialized)
	}
	s.log.Debug("Adding new votable option: %s", votableId)

	//checking if option already added
	optionExists, err := s.IsVotableOptionExists(ctx, votableId)
	if err != nil {
		return err
	}
	if optionExists {
		return Error(ErrVotingOptionAlreadyExists, votableId)
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

	s.log.Debug("Creating new asset for this votable id in voting ledger: %s", votableId)
	putStateErr := ctx.GetStub().PutState(VOTE_OPTION_REGISTER_PREFIX+votableId, votableOptionJSON) // new state added to the voting ledger
	return putStateErr

}

// Verify Signature
func VerifySignature(publicKey *rsa.PublicKey, data string, signature []byte) error {
	msgHash := sha512.New()
	_, err := msgHash.Write([]byte(data))
	if err != nil {
		return WrapError(ErrHashingData, err)
	}

	msgHashSum := msgHash.Sum(nil)

	err = rsa.VerifyPSS(publicKey, crypto.SHA512, msgHashSum, signature, nil)
	if err != nil {
		return WrapError(ErrSignatureValidation, err)
	}

	return nil
}

// function to add Voter
func (s *SmartContract) AddVoters(ctx contractapi.TransactionContextInterface, voterIds []string, signature_encrypted_base64 string) (error, []error) {

	if s.Initialized != true {
		return Error(ErrLedgerNotInitialized), []error{}
	}
	s.log.Debug("Decrypting signature ...")
	//Get the ledger's admin private key
	privateKey_bytes, err := base64.StdEncoding.DecodeString(s.PrivateKey)
	if err != nil {
		return WrapError(ErrDecodingBase64, err, "[SmartContract.PrivateKey]"), []error{}
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKey_bytes)
	if err != nil {
		return WrapError(ErrParsingPvtKey, err), []error{}
	}

	signature_encrypted, err := base64.StdEncoding.DecodeString(signature_encrypted_base64)
	if err != nil {
		return WrapError(ErrDecodingBase64, err, "[signature_encrypted_base64]"), []error{}
	}
	signature, err := DecryptOAEP(privateKey, signature_encrypted)
	if err != nil {
		return WrapError(ErrDecryption, err, "[signature_encrypted]"), []error{}
	}

	s.log.Debug("Verifying signature ...")

	//Get bytes from base64 public key from contract
	publicKey_bytes, err := base64.StdEncoding.DecodeString(s.PublicKey)
	if err != nil {
		return WrapError(ErrDecodingBase64, err, "[SmartContract.PublicKey]"), []error{}
	}

	//Now convert public key bytes into RSA public key
	publicKey, err := x509.ParsePKCS1PublicKey(publicKey_bytes)
	if err != nil {
		return WrapError(ErrParsingPubKey, err, "[SmartContract.PublicKey]"), []error{}
	}

	//Verify signatur using public key
	err = VerifySignature(publicKey, strings.Join(voterIds, ","), signature)
	if err != nil {
		return WrapError(ErrSignatureValidation, err), []error{}
	}

	errors := []error{}
	for _, voterId := range voterIds {
		s.log.Debug("Adding new voter: %s", voterId)
		//checking if voter already added
		exists, err := s.isVoterExists(ctx, voterId)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if exists {
			errors = append(errors, Error(ErrVoterAlreadyExists, voterId))
			continue
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
			errors = append(errors, WrapError(ErrJsonMarshalling, err))
			continue
		}

		s.log.Debug("Creating new asset for this voter %s", voterId)
		err = ctx.GetStub().PutState(VOTER_REGISTER_PREFIX+voterId, voterJSON) // new state added to the voter ledger
		if err != nil {
			errors = append(errors, WrapError(ErrSettingState, err))
		}
	}

	if len(errors) > 0 {
		return Error(ErrCouldAddAddAllVoters, (String(len(errors)) + "/" + String(len(voterIds)))), errors

	}
	return nil, []error{}

}

func String(n int) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

// function to authorize Voter, this returns a public key specific to this voter's encryption private key
func (s *SmartContract) AuthVoter(ctx contractapi.TransactionContextInterface, voterId string, publicKey_base64 string, signature_base64 string) (string, error) {

	if s.Initialized != true {
		return "", Error(ErrLedgerNotInitialized)
	}
	s.log.Debug("Authorizing voter: %s", voterId)

	//checking if ballot exists
	ballot, err := s.GetBallot(ctx, voterId)
	if err != nil {
		return "", WrapError(ErrRetrivingState, err, voterId)
	}
	if ballot == nil {
		return "", Error(ErrVoterNotExists, voterId)
	}

	// verify signature
	s.log.Debug("Verifying signature ...")

	//1. Get bytes from base64 public key
	publicKey_bytes, err := base64.StdEncoding.DecodeString(publicKey_base64)
	if err != nil {
		return "", WrapError(ErrDecodingBase64, err, "[public key]")
	}

	//2. Now convert public key bytes into RSA public key
	publicKey, err := x509.ParsePKCS1PublicKey(publicKey_bytes)
	if err != nil {
		return "", WrapError(ErrParsingPubKey, err)
	}

	//3. Decode signature
	signature_bytes, err := base64.StdEncoding.DecodeString(signature_base64)
	if err != nil {
		return "", WrapError(ErrDecodingBase64, err, "[signature]")
	}

	//4. Verify signatur using public key
	err = VerifySignature(publicKey, voterId, signature_bytes)
	if err != nil {
		return "", WrapError(ErrSignatureValidation, err)
	}

	// provision key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", WrapError(ErrKeyGeneration, err)
	}

	//1. Get bytes of privatekey
	privateKey_bytes := x509.MarshalPKCS1PrivateKey(privateKey)

	ballot.PubKey = publicKey_base64

	//Encrypt signature

	enc_publicKey := privateKey.PublicKey
	encryptedBytes, err := EncryptOAEP(&enc_publicKey, []byte(signature_base64))
	if err != nil {
		return "", WrapError(ErrEncryption, err, "[signature]")
	}

	encryptedBytes_base64 := base64.StdEncoding.EncodeToString(encryptedBytes)

	//Store base64 encoded of encrypted signature
	ballot.Signature = string(encryptedBytes_base64)

	s.log.Debug("Creating new asset for this voter %s", voterId)
	ballotJSON, err := json.Marshal(ballot)
	if err != nil {
		return "", WrapError(ErrJsonMarshalling, err, "[Ballot]")
	}

	err = ctx.GetStub().PutState(VOTER_REGISTER_PREFIX+voterId, ballotJSON) // updated state added to the voter ledger

	if err != nil {
		return "", WrapError(ErrSettingState, err, VOTER_REGISTER_PREFIX+voterId)
	}

	//Now store the private key
	err = ctx.GetStub().PutPrivateData(VOTER_REGISTER_PREFIX+voterId, "privateKey", privateKey_bytes)

	if err != nil {
		return "", WrapError(ErrSettingPrivate, err, VOTER_REGISTER_PREFIX+voterId)
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
		return nil, Error(ErrLedgerNotInitialized)
	}
	votableOptionJSON, err := ctx.GetStub().GetState(VOTE_OPTION_REGISTER_PREFIX + votableId)
	if err != nil {
		return nil, WrapError(ErrRetrivingState, err, VOTE_OPTION_REGISTER_PREFIX+votableId)
	}
	if votableOptionJSON == nil {
		return nil, Error(ErrNoStateExists, VOTE_OPTION_REGISTER_PREFIX+votableId)
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
		return nil, Error(ErrLedgerNotInitialized)
	}
	voterJSON, err := ctx.GetStub().GetState(VOTER_REGISTER_PREFIX + voterId)
	if err != nil {
		return nil, WrapError(ErrRetrivingState, err, VOTER_REGISTER_PREFIX+voterId)
	}
	if voterJSON == nil {
		return nil, Error(ErrNoStateExists, VOTER_REGISTER_PREFIX+voterId)
	}

	var ballot Ballot
	err = json.Unmarshal(voterJSON, &ballot)
	if err != nil {
		return nil, err
	}

	return &ballot, nil
}

// function to cast vote, the option must be encrypted using the user's private key and whole data must be signed using server private key
func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, voterId string, votableIds_base64 string, signature_base64 string) error {
	if s.Initialized != true {
		return Error(ErrLedgerNotInitialized)
	}

	//Check if user is authorized

	//get ballot
	ballot, err := s.GetBallot(ctx, voterId)
	if err != nil {
		return WrapError(ErrGetBallot, err)
	}
	if ballot == nil {
		return Error(ErrNoStateExists, "Ballot", voterId)
	}

	if ballot.PubKey == "" {
		return Error(ErrNotAuthorized, voterId)
	}

	if ballot.Casted {
		return Error(ErrAlreadyVoted, voterId)
	}

	//Verify signatur using public key
	publicKey_bytes, err := base64.StdEncoding.DecodeString(ballot.PubKey)
	if err != nil {
		return WrapError(ErrDecodingBase64, err, "[Ballot.PubKey]")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicKey_bytes)
	if err != nil {
		return WrapError(ErrParsingPubKey, err, "[Ballot.PubKey]")
	}

	signature_bytes, err := base64.StdEncoding.DecodeString(signature_base64)
	if err != nil {
		return WrapError(ErrDecodingBase64, err, "[Ballot.Signature]")
	}

	err = VerifySignature(publicKey, voterId+votableIds_base64, signature_bytes)
	if err != nil {
		return WrapError(ErrSignatureValidation, err)
	}

	//Get the ballot's private key
	privateKey_bytes, err := ctx.GetStub().GetPrivateData(VOTER_REGISTER_PREFIX+voterId, "privateKey")
	if err != nil {
		return WrapError(ErrRetrivingPrivate, err, VOTER_REGISTER_PREFIX+voterId, "privateKey")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKey_bytes)
	if err != nil {
		return WrapError(ErrParsingPvtKey, err)
	}

	votableIds_bytes, err := base64.StdEncoding.DecodeString(votableIds_base64)
	if err != nil {
		return WrapError(ErrDecodingBase64, err, "[VotableIds]")
	}
	votableIds_json, err := DecryptOAEP(privateKey, votableIds_bytes)
	if err != nil {
		return WrapError(ErrDecryption, err, "[VotableIds]")
	}

	var picks VoterPicks
	err = json.Unmarshal(votableIds_json, &picks)
	if err != nil {
		return WrapError(ErrJsonUnmarshalling, err, "[VotableIds]")
	}

	if len(picks.VotableIds) == 0 {
		if s.Abstainable {
			picks.VotableIds = append(picks.VotableIds, Abstained)
		} else {
			return Error(ErrNoVoteIsZero)
		}
	}
	if len(picks.VotableIds) > 1 && s.SingleChoice {
		return Error(ErrNoVoteIsMoreThanOne)
	}

	//Update the timestamp
	ballot.Timestamp = time.Now().UnixNano()
	if s.IsAnonymous {
		ballot.Picks = nil
	} else {
		ballot.Picks = picks.VotableIds
	}
	ballot.Casted = true

	for _, optionId := range picks.VotableIds {
		votableOption, err := s.ReadOption(ctx, optionId)
		if err != nil {
			return WrapError(ErrGetOption, err)
		}
		votableOption.Count = votableOption.Count + 1 //Vote

		votableOptionJSON, err := json.Marshal(votableOption)
		if err != nil {
			return WrapError(ErrJsonMarshalling, err, "[VotableOption]")
		}
		updatestate_err := ctx.GetStub().PutState(VOTE_OPTION_REGISTER_PREFIX+optionId, votableOptionJSON)
		if updatestate_err != nil {
			return WrapError(ErrSettingState, err, VOTE_OPTION_REGISTER_PREFIX+optionId)
		}
	}

	//Update Ballot
	ballotJSON, err := json.Marshal(ballot)
	if err != nil {
		return WrapError(ErrJsonMarshalling, err, "[Ballot]")
	}
	updatestate_err := ctx.GetStub().PutState(VOTER_REGISTER_PREFIX+voterId, ballotJSON)
	if updatestate_err != nil {
		return WrapError(ErrSettingState, err, VOTE_OPTION_REGISTER_PREFIX+voterId)
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
		return nil, Error(ErrNoIterator)
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
		return nil, Error(ErrNoIterator)
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
