package contract

import (
	"fmt"
	"strings"
)

// Error Wrapper
type wrapError struct {
	msg  string
	args string
}

func (err wrapError) Error() string {
	if len(strings.TrimSpace(err.args)) > 0 {
		return fmt.Sprintf("%s: %s", err.msg, err.args)
	}
	return err.msg
}

func (err wrapError) wrap(args string) error {
	return wrapError{msg: err.msg, args: args}
}

var (
	//Error related to State Database
	ErrRetrivingState = wrapError{msg: "Unable to get value from state database"}
	ErrSettingState   = wrapError{msg: "Unable to set value to state database"}
	ErrNoStateExists  = wrapError{msg: "No data exists in state database"}

	//Error related to Private Database
	ErrRetrivingPrivate = wrapError{msg: "Unable to get value from private database"}
	ErrSettingPrivate   = wrapError{msg: "Unable to set value to private database"}
	ErrNoPrivateExists  = wrapError{msg: "No data exists in private database"}

	//Chaincode related
	ErrCCAlreadyInitialized      = wrapError{msg: "Chaincode already initialized"}
	ErrLedgerNotInitialized      = wrapError{msg: "Ledger not initialized"}
	ErrVotingOptionAlreadyExists = wrapError{msg: "This votable option already exists"}
	ErrVotingOptionNotExists     = wrapError{msg: "This votable option does not exist"}
	ErrVoterAlreadyExists        = wrapError{msg: "This voter already exists"}
	ErrVoterNotExists            = wrapError{msg: "This voter does not exist"}
	ErrNoVoteIsZero              = wrapError{msg: "Number of options chosen to vote can not be zero"}
	ErrNoVoteIsMoreThanOne       = wrapError{msg: "Number of options chosen to vote can not be more than one"}
	ErrGetBallot                 = wrapError{msg: "Error getting ballot"}
	ErrNotAuthorized             = wrapError{msg: "This voter is not authorized to vote"}
	ErrAlreadyVoted              = wrapError{msg: "This voter is already casted the vote"}
	ErrGetOption                 = wrapError{msg: "Error getting option"}
	ErrNoIterator                = wrapError{msg: "No result iterator found!"}
	ErrCouldAddAddAllVoters      = wrapError{msg: "Failed to add some of the voter(s)"}

	//Crypto Errors
	ErrDecodingBase64      = wrapError{msg: "Error in decoded base64 encoded value"}
	ErrParsingPubKey       = wrapError{msg: "Error in parsing public key"}
	ErrParsingPvtKey       = wrapError{msg: "Error in parsing private key"}
	ErrHashingData         = wrapError{msg: "Error in hashing data"}
	ErrSignatureValidation = wrapError{msg: "Error in verifying signature"}
	ErrKeyGeneration       = wrapError{msg: "Error in generating key-pair"}
	ErrEncryption          = wrapError{msg: "Error in encrypting data"}
	ErrDecryption          = wrapError{msg: "Error in decrypting data"}

	//JSON Errors
	ErrJsonMarshalling   = wrapError{msg: "Error in marshalling data to json"}
	ErrJsonUnmarshalling = wrapError{msg: "Error in unmarshalling json to daya"}
)

//Utility

func convertArgsToString(args []string) string {
	return strings.Join(args, ",")
}

func Error(err wrapError, args ...string) error {
	if len(args) > 0 {
		return err.wrap(convertArgsToString(args))
	}
	return err
}

func WrapError(err wrapError, inner error, args ...string) error {
	if len(args) > 0 {
		return err.wrap(inner.Error() + ": " + convertArgsToString(args))
	}
	return err.wrap(inner.Error())
}
