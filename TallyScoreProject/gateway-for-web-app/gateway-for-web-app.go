package main

// todo before running this code: start the CA servers

// contains gateway code for defining API endpoints for business(user) registration, scoring mechanism for businesses

import(
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
	"context"
)

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


		// for connecting to cc:
		mspID="Tally"  // membership service provider identifier
		TallyScoreCCName="tallyscore"
		channelname="tally"
		users_common_path="/home/ubuntu/fabric/tally-network/fabric-ca-servers/tally/client/users/"   // mspPath= users_common_path + "PANofuserTestO/msp"
)


type registrationRequest struct{
	PAN string `json:"PAN" binding:"required"`
	Name string `json:"Name" binding:"required"`
	PhoneNo string `json:"PhoneNo" binding:"required"`
	Address string `json:"Address" binding:"required"`
	LicenseType string `json:"LicenseType" binding:"required"`
	Score string `json:"Score" binding:"required"`
}

type detailsStructure struct{
	PrivateKey string `json:"PrivateKey"`
	PublicKey string `json:"PublicKey"`
}

type UpdateValueRequest struct {
	PAN  string `json:"PAN" binding:"required"`
	ChangeVal string `json:"ChangeVal" binding:"required"`
}

// type PathCollection struct{
// 	mspPath string 
// }

func printUsage() {
	panic("Format to register user:\n" +
		"go run . <user_PAN> <name> <phoneNo> <address> <license_type>" + "\n")
}

// var mspPath string
// var certPath string
// var keyPath string


func main(){

	router:= gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "POST, PUT",
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

	timeout := 50 * time.Second
	timeoutGroup := router.Group("/", TimeoutMiddleware(timeout))
	{
		timeoutGroup.PUT("/TallyScoreProject/performRegistration",performRegistration)
		timeoutGroup.POST("/TallyScoreProject/increaseTallyScore", increaseTallyScore)
		timeoutGroup.POST("/TallyScoreProject/decreaseTallyScore", decreaseTallyScore)
	}

	router.Run("0.0.0.0:8080")

}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}


func performRegistration(c *gin.Context){

	var request registrationRequest
	c.BindJSON(&request)
	PAN:=request.PAN
	Name:=request.Name
	PhoneNo:=request.PhoneNo
	Address:=request.Address
	LicenseType:=request.LicenseType
	Score:= request.Score

	fmt.Printf("Initiating registration of user %s with starting score %s\n", Name, Score)
	password, err := registerUser(PAN, Name, PhoneNo, Address, LicenseType)
	if err!=nil{
		fmt.Printf("Error in step 1\n")
		fmt.Errorf("Error in the process of registration of user\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	fmt.Printf("Password: %s\n", password)

	fmt.Printf("Initial stage of registration successful! Initiating enrollment of user now.\n")
	// write code to call enrollUser() function
	detailsAsset, mspPath, err:= enrollUser(PAN, password)
	if err!= nil{
		fmt.Errorf("Error in enrollment stage\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	// the mspPath obtained as the second return value of enrollUser() function call can be used in creating a user asset using the tallyscore chaincode
	certPath:= mspPath + "/signcerts/cert.pem"
	keyPath:= mspPath + "/keystore"
	peer:= "tbchlfdevpeer01" 
	domain:= "tally.tallysolutions.com"
	peer_port:="7051"
	peerEndpoint:= peer + "." + domain + ":" + peer_port
	gatewayPeer:= peer + "." + domain
	tlsCertPath:= "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/" + domain + "/peers/" + peer + "/tls/ca.crt" 	


	// creating company asset
	client, gw:= connect(peerEndpoint, certPath, keyPath, tlsCertPath, gatewayPeer)
	contract:= getContract(gw, TallyScoreCCName)
	createCompanyAsset(contract, PAN) // PAN will be the licenseId(i.e. unique id of the business)
	gw.Close()
	client.Close()   


	fmt.Printf("mspPath: %s", mspPath)

	detailsAssetJSON, err := json.Marshal(detailsAsset)
    if err != nil {
		fmt.Errorf("Error in enrollment stage- error in conversion to JSON format\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
        return
    }

	fmt.Printf("Priv Key: %s \n", detailsAsset.PrivateKey)
	fmt.Printf("Public Key: %s \n", detailsAsset.PublicKey)
	fmt.Printf("REGISTRATION OF USER SUCCESSFUL!\n")
	c.Data(http.StatusOK, "application/json", detailsAssetJSON)

}


func increaseTallyScore(c *gin.Context){

	var request UpdateValueRequest
	c.BindJSON(&request)
	PAN:= request.PAN 
	incVal:= request.ChangeVal

	mspPath:= users_common_path + PAN + "/msp"
	certPath:= mspPath + "/signcerts/cert.pem"
	keyPath:= mspPath + "/keystore/"
	peer:= "tbchlfdevpeer01" 
	domain:= "tally.tallysolutions.com"
	peer_port:="7051"
	peerEndpoint:= peer + "." + domain + ":" + peer_port
	gatewayPeer:= peer + "." + domain
	tlsCertPath:= "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/" + domain + "/peers/" + peer + "/tls/ca.crt" 	

	// getting the contract
	client, gw:= connect(peerEndpoint, certPath, keyPath, tlsCertPath, gatewayPeer)
	contract:= getContract(gw, TallyScoreCCName)

	fmt.Printf("PAN: %s, IncreaseValue: %s\n", PAN, incVal)
	fmt.Printf("CertPath: %s \n", certPath)
	fmt.Printf("KeyPath: %s \n \n \n", keyPath)
	evaluatedAsset, err:= contract.SubmitTransaction("IncreaseScore", PAN, incVal)
	fmt.Printf("\n-------------> After SubmitTransaction: O/p= %s \n Error= %s \n", string(evaluatedAsset), err)
	gw.Close()
	client.Close()   // IMPLEMENT THIS IN A TRY-CATCH-FINALLY BLOCK  , OR USE DEFER(preferred)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluatedAsset)))

}

func decreaseTallyScore(c *gin.Context){

	var request UpdateValueRequest
	c.BindJSON(&request)
	PAN:= request.PAN 
	decVal:= request.ChangeVal

	mspPath:= users_common_path + PAN + "/msp"
	certPath:= mspPath + "/signcerts/cert.pem"
	keyPath:= mspPath + "/keystore/"
	peer:= "tbchlfdevpeer01" 
	domain:= "tally.tallysolutions.com"
	peer_port:="7051"
	peerEndpoint:= peer + "." + domain + ":" + peer_port
	gatewayPeer:= peer + "." + domain
	tlsCertPath:= "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/" + domain + "/peers/" + peer + "/tls/ca.crt" 	

	// getting the contract
	client, gw:= connect(peerEndpoint, certPath, keyPath, tlsCertPath, gatewayPeer)
	contract:= getContract(gw, TallyScoreCCName)

	fmt.Printf("PAN: %s, DecreaseValue: %s\n", PAN, decVal)
	evaluatedAsset, err:= contract.SubmitTransaction("DecreaseScore", PAN, decVal)
	fmt.Printf("\n-------------> After SubmitTransaction: O/p= %s \n Error= %s \n", string(evaluatedAsset), err)
	gw.Close()
	client.Close()   // IMPLEMENT THIS IN A TRY-CATCH-FINALLY BLOCK  , OR USE DEFER(preferred)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluatedAsset)))

}


func createCompanyAsset(contract *client.Contract, businessPAN string){
	fmt.Printf("After registration and enrollment, the score of the business will now be initialized.\n")
	result,err:= contract.SubmitTransaction("RegisterCompany", businessPAN)  // don't pass the PAN. 
	fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)
}


func getContract(gw *client.Gateway , ccName string) *client.Contract {
	network := gw.GetNetwork(channelname)
	return  network.GetContract(ccName)
}



func connect(peerEndpoint string, certPath string, keyPath string, tlsCertPath string, gatewayPeer string) (*grpc.ClientConn, *client.Gateway) {
	fmt.Printf("\nConnecting to : %s \n", peerEndpoint)

	// gRPC client conn- shared by all gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	//creating client identity, signing implementation
	id := newIdentity(certPath)          
	sign := newSign(keyPath)

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	return clientConnection, gw
}

func newGrpcConnection(tlsCertPath string, gatewayPeer string, peerEndpoint string) *grpc.ClientConn {

	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

func newIdentity(certPath string) *identity.X509Identity { 
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}
	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}
	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

func newSign(keyPath string) identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}


func registerUser(PAN string, Name string, PhoneNo string, Address string, LicenseType string) (string, error){   // this function should take in PAN and return the password

	cmdVariable := exec.Command("fabric-ca-client", "register", 
	"--id.name", PAN, 
	"--id.type", "client", 
	"--id.affiliation", "tally", 
	"--id.maxenrollments", "1", 
	"--id.attrs", fmt.Sprintf("pan=%s,name=%s,phone=%s,address=%s,license=%s", PAN, Name, PhoneNo, Address, LicenseType),
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

func enrollUser(PAN string, password string) (*detailsStructure, string, error) {  // this function takes in PAN and password, then it should returns the public+private key msp as a structure

	// urlmid would be like-> <PAN>:<password>

	mspPath := fmt.Sprintf("%s/users/%s", fabric_ca_client_home, PAN) +"/msp"
	cmdVariable:= exec.Command("fabric-ca-client", "enroll", "-u", urlstart + PAN + ":" + password + urlend , "--csr.names", "C=IN,ST=Karnataka,L=Bengaluru,O=Tally,OU=client", "-M", mspPath, "--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))
	fmt.Printf("%s", cmdVariable.String())
	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))
	err := cmdVariable.Run()
	fmt.Printf("%v \n", err)
	if err != nil {
		return nil, "", err
	}
	// return content of mspPath- with the signcert content(public key) as a param, private key as a param
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
		return nil, "", err
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

	// client side will be recreating private key and public key files content
	detailsAsset:= detailsStructure{
		PrivateKey: privatekey,
		PublicKey: publickey,
	}

	return &detailsAsset, mspPath, nil
}

func getPassword(outputString string) string{ // function to extract password from the output generated in the registerUser() function
	PasswordTextIndex := strings.Index(outputString, "Password: ")
	if PasswordTextIndex == -1 {
		return ""
	}
	password := outputString[PasswordTextIndex+len("Password: "):]
	return strings.TrimSpace(password)
}
