package main

import(
	"fmt"
	"github.com/hyperledger/fabric-ca-client-go/caclient"
	"github.com/hyperledger/fabric-ca-client-go/api"
)

func RegisterUser(caClient *caclient.Client, Userid string) (string, error) {

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

func EnrollUser(){
	
}
