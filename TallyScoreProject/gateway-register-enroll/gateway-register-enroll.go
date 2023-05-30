package main

import(
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// PRIOR TO RUNNING THIS CODE- START THE CA SERVERS: Navigate to Setup-Network and run ./2A_StartCAServer.sh 


// DELETE THIS COMMENT WHEN DONE WITH CODE: /home/ubuntu/fabric/tally-network/fabric-ca-servers/tally/client/users/PANofuserTestO/msp


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

type detailsStructure struct{
	PrivateKey string `json:"PrivateKey"`
	PublicKey string `json:"PublicKey"`
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
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	fmt.Printf("Password: %s\n", password)

	fmt.Printf("Initial stage of registration successful! Initiating enrollment of user now.\n")
	// write code to call enrollUser() function
	detailsAsset, err:= enrollUser(PAN, password)
	if err!= nil{
		fmt.Errorf("Error in enrollment stage\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}

	detailsAssetJSON, err := json.Marshal(detailsAsset)
    if err != nil {
		fmt.Errorf("Error in enrollment stage- error in conversion to JSON format\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
        return
    }

	fmt.Printf("Priv Key: %s \n", detailsAsset.PrivateKey)
	fmt.Printf("Public Key: %s \n", detailsAsset.PublicKey)
	fmt.Printf("REGISTRATION OF USER SUCCESSFUL!\n")

	// c.Writer.Header().Set("Content-Type","application/json")
	c.JSON(http.StatusOK, detailsAssetJSON)

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

func enrollUser(PAN string, password string) (*detailsStructure, error) {  // this function should take in PAN and password, then it should return/print the public+private key msp

	// urlmid would be like-> <PAN>:<password>

	mspPath := fmt.Sprintf("%s/users/%s", fabric_ca_client_home, PAN) +"/msp"
	cmdVariable:= exec.Command("fabric-ca-client", "enroll", "-u", urlstart + PAN + ":" + password + urlend , "--csr.names", "C=IN,ST=Karnataka,L=Bengaluru,O=Tally,OU=client", "-M", mspPath, "--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))
	fmt.Printf("%s", cmdVariable.String())
	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))
	err := cmdVariable.Run()
	fmt.Printf("%v \n", err)
	if err != nil {
		return nil, err
	}
	// return content of mspPath- with the signcert as a param, private key as a param, tls and ca-certs  then DELETE the mspPath folder

	// mspPath will have the path till the folder "msp"- which contains keystore(PRIVATE KEY location) and signcerts(containing cert.pem- from which the public key is to be extracted)
	// Extracting the private key
	pathKeystore:= mspPath + "/keystore"	// for private key
	pathSigncertFile:= mspPath + "/signcerts/cert.pem"	// for public key
	fmt.Printf("sign cert path: %s \n", pathSigncertFile)
	
	// below will be the default values
	privatekey:="private_key"
	
	
	files, err := ioutil.ReadDir(pathKeystore)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if !file.IsDir() {
			filePath := filepath.Join(pathKeystore, filename)

			// Read the contents of the file
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Println("Error reading file:", err)
				continue
			}

			// Process the private key content as per your requirements
			// In this example, we'll simply print it
			privatekey=string(content)
		}
	}

	// now to retrieve the public key
	certFileread, err := ioutil.ReadFile(pathSigncertFile)
	if err != nil {
		log.Fatal(err)
	}

	publickey:= string(certFileread)

	// decodedBlock, _ := pem.Decode(certFileread)  // decoding the cert file reaf
	// if decodedBlock == nil {
	// 	log.Fatal("Error in decoding PEM block")
	// }
	// cert, err := x509.ParseCertificate(decodedBlock.Bytes)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// publicKeyInterface := cert.PublicKey

	// publickey:= getPublicKeyAlgorithm(publicKeyInterface) + "\n" + getPublicKeyDetails(publicKeyInterface)

	// get the contents of cert.pem for public key and pass it to the client- in order to insert it into a cert.pem


	//now we create a new details asset that contains the public key and private key of the registering business

	// client side will be recreating private key and public key content
	detailsAsset:= detailsStructure{
		PrivateKey: privatekey,
		PublicKey: publickey,
	}

	return &detailsAsset,nil
}

func getPassword(outputString string) string{ // function to extract password from the output generated in the registerUser() function
	PasswordTextIndex := strings.Index(outputString, "Password: ")
	if PasswordTextIndex == -1 {
		return ""
	}
	password := outputString[PasswordTextIndex+len("Password: "):]
	return strings.TrimSpace(password)
}


func getPublicKeyAlgorithm(publicKeyInterface interface{}) string {
	switch publicKeyInterface.(type) {
	case *rsa.PublicKey:
		return "RSA"
	case *dsa.PublicKey:
		return "DSA"
	case *ecdsa.PublicKey:
		return "ECDSA"
	default:
		return "Unknown"
	}
}

func getPublicKeyDetails(publicKeyInterface interface{}) string {
	switch publicKeyInterface := publicKeyInterface.(type) {
	case *rsa.PublicKey:
		return fmt.Sprintf("RSA Key: Size %d bits", publicKeyInterface.Size()*8)
	case *dsa.PublicKey:
		return fmt.Sprintf("DSA Key: Size %d bits", publicKeyInterface.Q.BitLen())
	case *ecdsa.PublicKey:
		return fmt.Sprintf("ECDSA Key: Curve %s, Size %d bits", publicKeyInterface.Curve.Params().Name, publicKeyInterface.Curve.Params().BitSize)
	default:
		return "Unknown"
	}

}
