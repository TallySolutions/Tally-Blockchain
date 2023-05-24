package main

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// PRIOR TO RUNNING THIS CODE- START THE CA SERVERS: Navigate to Setup-Network and run ./2A_StartCAServer.sh 


var(
	tallyHome string
	caServerHome string
	tallyCAHome string
	fabric_ca_client_home string
	urlend string
)

const(
		networkHome= "fabric/tally-network"
		tallyCAName= "tally"
		ca_host="tbchlfdevca01"
		domain="tally.tallysolutions.com"
		tally_ca_port="7055"
		urlstart= "https://"
)


type registrationRequest struct{
	PAN string `json:"PAN" binding:"required"`
	Name string `json:"Name" binding:"required"`
	PhoneNo string `json:"PhoneNo" binding:"required"`
	Address string `json:"Address" binding:"required"`
	LicenseType string `json:"LicenseType" binding:"required"`
}

func printUsage() {
	panic("Format to register user:\n" +
		"go run . <user_PAN> <name> <phoneNo> <address> <license_type>" + "\n")
}

func main(){


	router:= gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "POST",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge: 50 * time.Second,
		Credentials: false,
		ValidateHeaders: false,
	}))

	tallyHome= os.Getenv("HOME") + "/" + networkHome
	caServerHome= tallyHome + "/fabric-ca-servers"
	tallyCAHome= caServerHome + "/" + tallyCAName
	fabric_ca_client_home= tallyCAHome + "/client"
	urlend= "@" + ca_host + "." + domain + ":" + tally_ca_port

	router.POST("/TallyScoreProject/performRegistration",performRegistration)
	router.Run("0.0.0.0:8080")

}

func performRegistration(c *gin.Context){

	var request registrationRequest
	c.BindJSON(&request)
	PAN:=request.PAN
	name:=request.Name
	phoneNo:=request.PhoneNo
	address:=request.Address
	license:=request.LicenseType

	fmt.Printf("Initiating registration of user %s\n", name)
	password, err := registerUser(PAN, name, phoneNo, address, license)
	if err!=nil{
		fmt.Printf("Error in step 1\n")
		fmt.Errorf("Error in the process of registration of user\n")
		return
	}
	fmt.Printf("Password: %s\n", password)

	fmt.Printf("Initial stage of registration successful! Initiating enrollment of user now.\n")
	// write code to call enrollUser() function
	userMSP, err:= enrollUser(PAN, password)
	if err!= nil{
		fmt.Errorf("Error in enrollment stage\n")
		return
	}
	fmt.Printf("MSP path of User: %s \n", userMSP)
	fmt.Printf("Registration of User successful!\n")

	c.Writer.Header().Set("Content-Type","application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", userMSP))


}


func registerUser(PAN string, name string, phoneNo string, address string, license string) (string, error){   // this function should take in PAN and print the password

	cmdVariable := exec.Command("fabric-ca-client", "register", 
	"--id.name", PAN, 
	"--id.type", "client", 
	"--id.affiliation", "tally", 
	"--id.maxenrollments", "1", 
	"--id.attrs", fmt.Sprintf("pan=%s,name=%s,phone=%s,address=%s,license=%s", PAN, name, phoneNo, address, license),
	"--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))


	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))

	fmt.Printf("cmd Env: %s\n",cmdVariable.Env)

	fmt.Printf("FABRIC CA CLIENT HOME PATH: %s \n", fabric_ca_client_home)


	fmt.Printf("The command executed while registering:%s\n", cmdVariable.String())
	fmt.Print("Value of tallyCAHome:", tallyCAHome, "\n")
	fmt.Printf("Path to ca-certs:%s \n", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))


	output, err := cmdVariable.CombinedOutput()
	if err != nil {
		return "",err
	}
	
	password := getPassword(string(output)) // extract password from the cli's output
	return password,nil

}

func enrollUser(PAN string, password string) (string, error) {  // this function should take in PAN and password, then it should return/print the public+private key msp

	// urlmid would be like-> <PAN>:<password>

	mspPath := fmt.Sprintf("%s/users/%s", fabric_ca_client_home, PAN) +"/msp"
	cmdVariable:= exec.Command("fabric-ca-client", "enroll", "-u", urlstart + PAN + ":" + password + urlend , "--csr.names", "C=IN,ST=Karnataka,L=Bengaluru,O=Tally,OU=client", "-M", mspPath, "--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))
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