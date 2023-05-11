package main



// TRY TO USE OS.EXEC TO CALL FABRIC-CA CLI COMMANDS-- or look into forking fabric ca client

import(
	"github.com/hyperledger/fabric"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-ca/api"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	cspsigner "github.com/hyperledger/fabric/bccsp/signer"
	"github.com/hyperledger/fabric/bccsp/utils"

)

func RegisterUser(caClient *client.Client, Userid string) (string, error) {

	RegisterReq := &api.RegistrationRequest{  // registration request is initiated
		Name:   Userid,
		Secret: "tally-password",
		Type:   "client",
		MaxEnrollments: 1,
	}

	password, err := caClient.Register(RegisterReq)   // registration of the user on the server, password for that user is returned 
	if err != nil {
		return "", fmt.Errorf("error in registering user: %v", err)
	}

	return password, nil   // this password is to be used for enrollment
}

func EnrollUser(caClient *client.Client, Userid string, password string){
	
	EnrollReq := &api.EnrollmentRequest{
        Name:     Userid,
        Secret:   password,
        Profile:  "tls",
        Label:    "",
        CSRHosts: nil,
		CSRNames: []csr.Name{
			{
				C:  "IN",
				ST: "Bengaluru",
				L:  "Bengaluru",
				O:  "Tally",
				OU: "client",
			},
   	 	}
	}

	enrollResponse, err := caClient.Enroll(EnrollReq)
    if err != nil {
        return nil, fmt.Errorf("failed to enroll user: %v", err)
    }

    identity, err := identity.NewX509Identity(userID, enrollResponse.Cert, enrollResponse.Key)
    if err != nil {
        return nil, fmt.Errorf("failed to create identity: %v", err)
    }

    return identity, nil
}
