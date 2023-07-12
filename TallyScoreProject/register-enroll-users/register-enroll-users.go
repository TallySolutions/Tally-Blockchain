package main

import(
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var tallyHome string
var caServerHome string
var tallyCAHome string
var fabric_ca_client_home string
var urlend string

const(
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


func main(){

	tallyHome= os.Getenv("HOME") + "/" + networkHome
	caServerHome= tallyHome + "/fabric-ca-servers"
	tallyCAHome= caServerHome + "/" + tallyCAName
	fabric_ca_client_home= tallyCAHome + "/client"
	urlend= "@" + ca_host + "." + domain + ":" + tally_ca_port

		userId:= os.Args[1]
		fmt.Printf("Initiating registration of user %s\n", userId)
		password, err := registerUser(userId)
		if err!=nil{
			fmt.Printf("Error in step 1\n")
			fmt.Errorf("Error in the process of registration of user\n")
			return
		}
		fmt.Printf("Password: %s\n", password)

		fmt.Printf("Initial stage of registration successful! Initiating enrollment of user now.\n")
		// write code to call enrollUser() function
		userMSP, err:= enrollUser(userId, password)
		if err!= nil{
			fmt.Errorf("Error in enrollment stage\n")
			return
		}
		fmt.Printf("MSP path of User: %s \n", userMSP)
		fmt.Printf("Registration of User successful!\n")

}


// func performRegistraion(userId)

func registerUser(userId string) (string, error){   // this function should take in userid and print the password

	cmdVariable := exec.Command("fabric-ca-client", "register", "--id.name", userId, "--id.type", "client", "--id.affiliation", "tally", "--id.maxenrollments", "1", "--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))
	// set max enrollments
	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))

	fmt.Printf("cmd Env: %s\n",cmdVariable.Env)

	fmt.Printf("FABRIC CA CLIENT HOME PATH: %s \n", fabric_ca_client_home)


	fmt.Printf("The command executed while registering:%s\n", cmdVariable.String())
	fmt.Print("Value of tallyCAHome:", tallyCAHome, "\n")
	fmt.Printf("Path to ca-certs:%s \n", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))


	output, err := cmdVariable.CombinedOutput()

	fmt.Printf("Value of cmdVariable's output string: %s\n", string(output) )

	if err != nil {
		return "",err
	}
	
	password := getPassword(string(output)) // extract password from the cli's output
	return password,nil

}

func enrollUser(userId string, password string) (string, error) {  // this function should take in userid and password, then it should return/print the public+private key msp

	// urlmid would be like-> <userId>:<password>

	mspPath := fmt.Sprintf("%s/users/%s", fabric_ca_client_home, userId) +"/msp"
	cmdVariable:= exec.Command("fabric-ca-client", "enroll", "-u", urlstart + userId + ":" + password + urlend , "--csr.names", "C=IN,ST=Karnataka,L=Bengaluru,O=Tally,OU=client", "-M", mspPath, "--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))
	fmt.Printf("%s", cmdVariable.String())
	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))
	err := cmdVariable.Run()
	fmt.Printf("%v \n", err)
	if err != nil {
		return "", err
	}
	// return content of mspPath- with the signcert as a param, private key as a param, tls and ca-certs  then DELETE the mspPath folder
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