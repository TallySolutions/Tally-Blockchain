package main

import(
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// make sure to start the fabric ca servers- before running this


var tallyHome string
var caServerHome string
var tallyCAHome string
var fabric_ca_client_home string
var urlend string

const(
		// add the relevant env vars(from setup network) as constants
		// var- TALLY NETWORK HOME- add the env locations with this as the common base
		networkHome= "fabric/tally-network"
		tallyCAName= "tally"
		ca_host="tbchlfdevca01"
		domain="tally.tallysolutions.com"
		tally_ca_port="7055"
		urlstart= "https://"
)


func printUsage() {
	panic("Format to register user:\n" +
		"go run . <user_id>" + "\n")
}

// func printUsage() {
// 	panic("Format to deal with users:\n" +
// 		"Register <user_id>\n" +
// 		"Enroll <user_id> <password_generated_after_registering>" + "\n")
// }

func main(){

	tallyHome= os.Getenv("HOME") + "/" + networkHome
	caServerHome= tallyHome + "/fabric-ca-servers"
	tallyCAHome= caServerHome + "/" + tallyCAName
	fabric_ca_client_home= tallyCAHome + "/client"
	urlend= "@" + ca_host + "." + domain + ":" + tally_ca_port

	if len(os.Args)<2{
		printUsage()
	}

	operation:= os.Args[0]
	if operation == "Register" && len(os.Args)<2{
		printUsage()
	}

	// if operation == "Register"{
		userId:=os.Args[1]
		fmt.Printf("Initiating registration of user\n")
		password, err := registerUser(userId)
		if err!=nil{
			fmt.Errorf("Error in the process of registration of user\n")
		}
		fmt.Printf("Password: %s\n", password)

		fmt.Printf("Initial stage of registration successful! Initiating enrollment of user now.\n")
		// write code to call enrollUser() function
		userMSP, err:= enrollUser(userId, password)
		if err!= nil{
			fmt.Errorf("Error in enrollment stage\n")
		}
		fmt.Printf("MSP path of User: %s \n", userMSP)
		fmt.Printf("Registration of User successful!\n")

//	}

}

func registerUser(userId string) (string, error){   // this function should take in userid and print the password

	cmdVariable := exec.Command("fabric-ca-client", "register", "--id.name", userId, "--id.type", "client", "--id.affiliation", "tally", "--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))
	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))
	output, err := cmdVariable.Output()
	if err != nil {
		return "",err
	}
	password := getPassword(string(output)) // extract password from the cli's output
	fmt.Printf("Password:%s", password)
	return password,nil

}

func enrollUser(userId string, password string) (string, error) {  // this function should take in userid and password, then it should return/print the public+private key msp

	// urlmid would be like-> <userId>:<password>
	cmdVariable:= exec.Command("fabric-ca-client", "enroll", "-u", urlstart + userId + ":" + password + urlend)
	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))
	err := cmdVariable.Run()
	if err != nil {
		return "", err
	}
	
	mspPath := fmt.Sprintf("%s/msp", fabric_ca_client_home)

	return mspPath, nil

}

func getPassword(outputString string) string{ // function to extract password from the output generated in the registerUser() function
	PasswordTextIndex := strings.Index(outputString, "Password: ")
	if PasswordTextIndex == -1 {
		return ""
	}
	password := outputString[PasswordTextIndex+len("Password: "):]
	return strings.TrimSpace(password)
}