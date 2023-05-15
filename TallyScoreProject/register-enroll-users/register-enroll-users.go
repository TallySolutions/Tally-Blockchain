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

const(
		// add the relevant env vars(from setup network) as constants
		// var- TALLY NETWORK HOME- add the env locations with this as the common base
		networkHome= "fabric/tally-network"
		tallyCAName= "tally"
)

func printUsage() {
	panic("Format to deal with users:\n" +
		"Register <user_id>\n" +
		"Enroll <user_id> <password_generated_after_registering>" + "\n")
}

func main(){

	tallyHome= os.Getenv("HOME") + "/" + networkHome
	caServerHome= tallyHome + "/fabric-ca-servers"
	tallyCAHome= caServerHome + "/" + tallyCAName
	fabric_ca_client_home= tallyCAHome + "/client"

	if len(os.Args)<2{
		printUsage()
	}

	operation:=os.Args[0]
	if operation == "Register" && len(os.Args)<2{
		printUsage()
	}
	if operation == "Enroll" && len(os.Args)<3{
		printUsage()
	}

	if operation == "Register"{
		fmt.Printf("Intiating registration of user...\n")
		userId:=os.Args[1]
		registerUser(userId)
	} else if operation == "Enroll"{
		fmt.Printf("Intiating enrollment of user...\n")
		userId:=os.Args[1]
		password:=os.Args[2]
		enrollUser(userId, password)
	}

}

func registerUser(userId string) error{   // this function should take in userid and print the password
	cmdVariable := exec.Command("fabric-ca-client", "register", "--id.name", userId, "--id.type", "client", "--id.affiliation", "tally", "--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))
	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))
	output, err := cmdVariable.Output()
	if err != nil {
		return err
	}

	password := getPassword(string(output)) // extract password from the cli's output
	fmt.Printf("Password:%s", password)
	return nil
	

}

func enrollUser(userId string, password string){  // this function should take in userid and password, then it should return/print the public+private key msp
	// USE ID.SECRETS FOR PASSWORD- do not make this a constant in the prog

}

func getPassword(outputString string) string{ // function to extract password from the output generated in the registerUser() function
	PasswordTextIndex := strings.Index(outputString, "Password: ")
	if PasswordTextIndex == -1 {
		return ""
	}
	password := outputString[PasswordTextIndex+len("Password: "):]
	return strings.TrimSpace(password)
}