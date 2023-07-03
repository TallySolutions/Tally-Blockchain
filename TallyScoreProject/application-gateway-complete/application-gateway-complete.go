package main

// todo before running this code: start the CA servers

// contains gateway code for defining API endpoints for business(user) registration, scoring mechanism for businesses, all voucher related mechanisms

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	cors "github.com/itsjamie/gin-cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tallyHome             string
	caServerHome          string
	tallyCAHome           string
	fabric_ca_client_home string
	urlend                string
)

const (
	networkHome   = "fabric/tally-network"
	tallyCAName   = "tally"
	ca_host       = "tbchlfdevca01"
	domain        = "tally.tallysolutions.com"
	tally_ca_port = "7055"
	urlstart      = "https://"

	// for connecting to cc:
	mspID                 = "Tally" // membership service provider identifier
	TallyScoreCCName      = "tallyscore"
	BusinessProfileCCName = "businessprofile"
	channelname           = "tally"
	users_common_path     = "/home/ubuntu/fabric/tally-network/fabric-ca-servers/tally/client/users/" // mspPath= users_common_path + "PANofuserTestO/msp"
)

type registrationRequest struct {
	PAN         string `json:"PAN" binding:"required"`
	Name        string `json:"Name" binding:"required"`
	PhoneNo     string `json:"PhoneNo" binding:"required"`
	Address     string `json:"Address" binding:"required"`
	LicenseType string `json:"LicenseType" binding:"required"`
	Score       string `json:"Score" binding:"required"`
}

type detailsStructure struct {
	PrivateKey string `json:"PrivateKey"`
	PublicKey  string `json:"PublicKey"`
}

type UpdateValueRequest struct {
	PAN       string `json:"PAN" binding:"required"`
	ChangeVal string `json:"ChangeVal" binding:"required"`
}

type voucherCreationRequest struct {
	VoucherID   string `json:"VoucherID" binding:"required"`
	SupplierID  string `json:"SupplierID" binding:"required"`
	VoucherType string `json:"VoucherType" binding:"required"`
	Hashcode    string `json:"Hashcode" binding:"required"`
	TotalValue  string `json:"TotalValue" binding:"required"`
	Currency    string `json:"Currency" binding:"required"`
}

type voucherIDRequest struct {
	VoucherID string `json:"VoucherID" binding:"required"`
}

type voucherUpdationRequest struct {
	VoucherID    string `json:"VoucherID" binding:"required"`
	Parameter    string `json:"Parameter" binding:"required"`
	UpdatedValue string `json:"UpdatedValue" binding:"required"`
}

func main() {

	router := gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, POST, PUT",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	tallyHome = os.Getenv("HOME") + "/" + networkHome
	caServerHome = tallyHome + "/fabric-ca-servers"
	tallyCAHome = caServerHome + "/" + tallyCAName
	fabric_ca_client_home = tallyCAHome + "/client"
	urlend = "@" + ca_host + "." + domain + ":" + tally_ca_port

	timeout := 50 * time.Second
	timeoutGroup := router.Group("/", TimeoutMiddleware(timeout))
	{
		timeoutGroup.PUT("/TallyScoreProject/performRegistration", performRegistration) // to register + enroll the business
		timeoutGroup.POST("/TallyScoreProject/increaseTallyScore", increaseTallyScore)  // to increae tallyscore
		timeoutGroup.POST("/TallyScoreProject/decreaseTallyScore", decreaseTallyScore)  // to decrease tallyscore

		timeoutGroup.PUT("/TallyScoreProject/voucherCreation/:PAN", voucherCreation)
		timeoutGroup.POST("/TallyScoreProject/voucherCancellation/:PAN", voucherCancellation)
		timeoutGroup.PUT("/TallyScoreProject/voucherUpdation/:PAN", voucherUpdation)
		timeoutGroup.GET("/TallyScoreProject/listOwnerVouchers/:PAN", listOwnerVouchers) // above 4 are owner related voucher endpoints

		timeoutGroup.GET("/TallyScoreProject/readVoucher/:PAN", readVoucher)

		timeoutGroup.POST("/TallyScoreProject/voucherApproval/:PAN", voucherApproval)
		timeoutGroup.POST("/TallyScoreProject/voucherReturn/:PAN", voucherReturn)
		timeoutGroup.POST("/TallyScoreProject/voucherRejection/:PAN", voucherRejection)
		timeoutGroup.GET("/TallyScoreProject/listSupplierVouchers/:PAN", listSupplierVouchers) // above 4 are supplier related voucher endpoints
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

func performRegistration(c *gin.Context) {

	var request registrationRequest
	c.BindJSON(&request)
	PAN := request.PAN
	Name := request.Name
	PhoneNo := request.PhoneNo
	Address := request.Address
	LicenseType := request.LicenseType
	Score := request.Score

	fmt.Printf("Initiating registration of user %s with starting score %s\n", Name, Score)
	password, err := registerUser(PAN, Name, PhoneNo, Address, LicenseType)
	if err != nil {
		fmt.Printf("Error in step 1\n")
		// fmt.Errorf("Error in the process of registration of user\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	fmt.Printf("Password: %s\n", password)

	fmt.Printf("Initial stage of registration successful! Initiating enrollment of user now.\n")
	// write code to call enrollUser() function
	detailsAsset, mspPath, err := enrollUser(PAN, password)
	if err != nil {
		// fmt.Errorf("Error in enrollment stage\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// the mspPath obtained as the second return value of enrollUser() function call can be used in creating a user asset using the tallyscore chaincode
	certPath := mspPath + "/signcerts/cert.pem"
	keyPath := mspPath + "/keystore"
	peer := "tbchlfdevpeer01"
	domain := "tally.tallysolutions.com"
	peer_port := "7051"
	peerEndpoint := peer + "." + domain + ":" + peer_port
	gatewayPeer := peer + "." + domain
	tlsCertPath := "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/" + domain + "/peers/" + peer + "/tls/ca.crt"

	// creating company asset
	client, gw := connect(peerEndpoint, certPath, keyPath, tlsCertPath, gatewayPeer)
	contract := getContract(gw, TallyScoreCCName)
	createCompanyAsset(contract, PAN) // PAN will be the licenseId(i.e. unique id of the business)
	gw.Close()
	client.Close()

	fmt.Printf("mspPath: %s", mspPath)

	detailsAssetJSON, err := json.Marshal(detailsAsset)
	if err != nil {
		// fmt.Errorf("Error in enrollment stage- error in conversion to JSON format\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	fmt.Printf("Priv Key: %s \n", detailsAsset.PrivateKey)
	fmt.Printf("Public Key: %s \n", detailsAsset.PublicKey)
	fmt.Printf("REGISTRATION OF USER SUCCESSFUL!\n")
	c.Data(http.StatusOK, "application/json", detailsAssetJSON)

}

func increaseTallyScore(c *gin.Context) {

	var request UpdateValueRequest
	c.BindJSON(&request)
	PAN := request.PAN
	incVal := request.ChangeVal

	mspPath := users_common_path + PAN + "/msp"
	certPath := mspPath + "/signcerts/cert.pem"
	keyPath := mspPath + "/keystore/"
	peer := "tbchlfdevpeer01"
	domain := "tally.tallysolutions.com"
	peer_port := "7051"
	peerEndpoint := peer + "." + domain + ":" + peer_port
	gatewayPeer := peer + "." + domain
	tlsCertPath := "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/" + domain + "/peers/" + peer + "/tls/ca.crt"

	// getting the contract
	client, gw := connect(peerEndpoint, certPath, keyPath, tlsCertPath, gatewayPeer)
	contract := getContract(gw, TallyScoreCCName)

	fmt.Printf("PAN: %s, IncreaseValue: %s\n", PAN, incVal)
	fmt.Printf("CertPath: %s \n", certPath)
	fmt.Printf("KeyPath: %s \n \n \n", keyPath)
	evaluatedAsset, err := contract.SubmitTransaction("IncreaseScore", PAN, incVal)
	fmt.Printf("\n-------------> After SubmitTransaction: O/p= %s \n Error= %s \n", string(evaluatedAsset), err)
	gw.Close()
	client.Close() // IMPLEMENT THIS IN A TRY-CATCH-FINALLY BLOCK  , OR USE DEFER(preferred)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluatedAsset)))

}

func decreaseTallyScore(c *gin.Context) {

	var request UpdateValueRequest
	c.BindJSON(&request)
	PAN := request.PAN
	decVal := request.ChangeVal

	mspPath := users_common_path + PAN + "/msp"
	certPath := mspPath + "/signcerts/cert.pem"
	keyPath := mspPath + "/keystore/"
	peer := "tbchlfdevpeer01"
	domain := "tally.tallysolutions.com"
	peer_port := "7051"
	peerEndpoint := peer + "." + domain + ":" + peer_port
	gatewayPeer := peer + "." + domain
	tlsCertPath := "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/" + domain + "/peers/" + peer + "/tls/ca.crt"

	// getting the contract
	client, gw := connect(peerEndpoint, certPath, keyPath, tlsCertPath, gatewayPeer)
	contract := getContract(gw, TallyScoreCCName)

	fmt.Printf("PAN: %s, DecreaseValue: %s\n", PAN, decVal)
	evaluatedAsset, err := contract.SubmitTransaction("DecreaseScore", PAN, decVal)
	fmt.Printf("\n-------------> After SubmitTransaction: O/p= %s \n Error= %s \n", string(evaluatedAsset), err)
	gw.Close()
	client.Close() // IMPLEMENT THIS IN A TRY-CATCH-FINALLY BLOCK  , OR USE DEFER(preferred)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluatedAsset)))

}

func createCompanyAsset(contract *client.Contract, businessPAN string) {
	fmt.Printf("After registration and enrollment, the score of the business will now be initialized.\n")
	result, err := contract.SubmitTransaction("RegisterCompany", businessPAN) // don't pass the PAN.
	fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)
}

func getContract(gw *client.Gateway, ccName string) *client.Contract {
	network := gw.GetNetwork(channelname)
	return network.GetContract(ccName)
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

func registerUser(PAN string, Name string, PhoneNo string, Address string, LicenseType string) (string, error) { // this function should take in PAN and return the password

	cmdVariable := exec.Command("fabric-ca-client", "register",
		"--id.name", PAN,
		"--id.type", "client",
		"--id.affiliation", "tally",
		"--id.maxenrollments", "1",
		"--id.attrs", fmt.Sprintf("pan=%s,name=%s,phone=%s,address=%s,license=%s", PAN, Name, PhoneNo, Address, LicenseType),
		"--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))

	cmdVariable.Env = append(cmdVariable.Env, fmt.Sprintf("FABRIC_CA_CLIENT_HOME=%s", fabric_ca_client_home))

	fmt.Printf("cmd Env: %s\n", cmdVariable.Env)

	fmt.Printf("FABRIC CA CLIENT HOME PATH: %s \n", fabric_ca_client_home)

	fmt.Printf("The command executed while registering:%s\n", cmdVariable.String())
	fmt.Print("Value of tallyCAHome:", tallyCAHome, "\n")
	fmt.Printf("Path to ca-certs:%s \n", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))

	output, err := cmdVariable.CombinedOutput()
	if err != nil {
		return "", err
	}

	password := getPassword(string(output)) // extract password from the cli's output
	return password, nil

}

func enrollUser(PAN string, password string) (*detailsStructure, string, error) { // this function takes in PAN and password, then it should returns the public+private key msp as a structure

	// urlmid would be like-> <PAN>:<password>

	mspPath := fmt.Sprintf("%s/users/%s", fabric_ca_client_home, PAN) + "/msp"
	cmdVariable := exec.Command("fabric-ca-client", "enroll", "-u", urlstart+PAN+":"+password+urlend, "--csr.names", "C=IN,ST=Karnataka,L=Bengaluru,O=Tally,OU=client", "-M", mspPath, "--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", tallyCAHome))
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
	pathKeystore := mspPath + "/keystore"               // for private key
	pathSigncertFile := mspPath + "/signcerts/cert.pem" // for public key
	fmt.Printf("sign cert path: %s \n", pathSigncertFile)

	// below will be the default values
	privatekey := "private_key"

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

			privatekey = string(content)
		}
	}

	// Now, to retrieve the public key
	certFileread, err := ioutil.ReadFile(pathSigncertFile)
	if err != nil {
		log.Fatal(err)
	}

	publickey := string(certFileread)

	// Client side will be recreating private key and public key files content
	detailsAsset := detailsStructure{
		PrivateKey: privatekey,
		PublicKey:  publickey,
	}

	return &detailsAsset, mspPath, nil
}

func getPassword(outputString string) string { // function to extract password from the output generated in the registerUser() function
	PasswordTextIndex := strings.Index(outputString, "Password: ")
	if PasswordTextIndex == -1 {
		return ""
	}
	password := outputString[PasswordTextIndex+len("Password: "):]
	return strings.TrimSpace(password)
}

// Below are all the voucher related functions

func readVoucher(c *gin.Context) {

	var request voucherIDRequest
	VoucherID := c.Query("voucherID")
	request.VoucherID = VoucherID
	PAN := c.Param("PAN")

	fmt.Printf("Initiating reading of Voucher %s", VoucherID)
	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)

	asset, err := contract.EvaluateTransaction("ReadVoucher", VoucherID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Here is the string of asset read: %s \n", string(asset))
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(asset)))
	gw.Close()
	client.Close()

}

func voucherCreation(c *gin.Context) {

	var request voucherCreationRequest
	c.BindJSON(&request)
	VoucherID := request.VoucherID
	SupplierID := request.SupplierID
	VoucherType := request.VoucherType
	Hashcode := request.Hashcode
	TotalValue := request.TotalValue
	Currency := request.Currency

	PAN := c.Param("PAN")

	fmt.Printf("Initiating creation of voucher %s", VoucherID)

	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)
	fmt.Printf("Calling the contract named: %s \n", BusinessProfileCCName)
	result, err := contract.SubmitTransaction("VoucherCreated", VoucherID, SupplierID, VoucherType, Hashcode, TotalValue, Currency)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)

	asset, err := contract.SubmitTransaction("ReadVoucher", VoucherID)
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "application/json", asset)

	gw.Close()
	client.Close()

	// c.JSON(http.StatusOK, gin.H{"message": "Voucher created successfully"})
}

func voucherCancellation(c *gin.Context) {

	var request voucherIDRequest
	c.BindJSON(&request)
	VoucherID := request.VoucherID
	PAN := c.Param("PAN")

	fmt.Printf("Initiating cancellation of Voucher %s", VoucherID)
	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)

	result, err := contract.SubmitTransaction("VoucherCancelled", VoucherID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)
		return
	}
	fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)

	asset, err := contract.SubmitTransaction("ReadVoucher", VoucherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		panic(err)
	}
	c.Writer.Header().Set("Content-Type","application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(asset)))
	gw.Close()
	client.Close()

}

func voucherUpdation(c *gin.Context) {

	var request voucherUpdationRequest
	c.BindJSON(&request)
	VoucherID := request.VoucherID
	Parameter := request.Parameter
	UpdatedValue := request.UpdatedValue

	PAN := c.Param("PAN")

	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)

	asset, err := contract.SubmitTransaction("ReadVoucher", VoucherID)
	if err != nil {
		panic(err)
	}

	result, err := contract.SubmitTransaction("VoucherUpdated", VoucherID, Parameter, UpdatedValue)
	if err != nil {
		fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)
		return
	}

	asset, err = contract.SubmitTransaction("ReadVoucher", VoucherID)
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "application/json", asset)
	c.JSON(http.StatusOK, gin.H{"message": "Voucher updated successfully"})

	gw.Close()
	client.Close()
}

func listOwnerVouchers(c *gin.Context) {

	PAN := c.Param("PAN")
	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)
	assetsList, err := contract.SubmitTransaction("GetOwnerVouchers")
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "application/json", assetsList)

	gw.Close()
	client.Close()

}

func voucherApproval(c *gin.Context) {

	var request voucherIDRequest
	c.BindJSON(&request)
	VoucherID := request.VoucherID
	PAN := c.Param("PAN")

	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)
	result, err := contract.SubmitTransaction("VoucherApproved", VoucherID)

	if err != nil {
		fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	asset, err := contract.SubmitTransaction("ReadVoucher", VoucherID)
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "application/json", asset)
	gw.Close()
	client.Close()
}

func voucherRejection(c *gin.Context) {

	var request voucherIDRequest
	c.BindJSON(&request)
	VoucherID := request.VoucherID
	PAN := c.Param("PAN")

	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)
	result, err := contract.SubmitTransaction("VoucherRejected", VoucherID)

	if err != nil {
		fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	asset, err := contract.SubmitTransaction("ReadVoucher", VoucherID)
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "application/json", asset)
	c.JSON(http.StatusOK, gin.H{"message": "Voucher rejected successfully"})
	gw.Close()
	client.Close()
}

func voucherReturn(c *gin.Context) {

	var request voucherIDRequest
	c.BindJSON(&request)
	VoucherID := request.VoucherID
	PAN := c.Param("PAN")

	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)
	result, err := contract.SubmitTransaction("VoucherSentBack", VoucherID)

	if err != nil {
		fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	asset, err := contract.SubmitTransaction("ReadVoucher", VoucherID)
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "application/json", asset)
	c.JSON(http.StatusOK, gin.H{"message": "Voucher sent back successfully"})
	gw.Close()
	client.Close()
}

func listSupplierVouchers(c *gin.Context) {

	PAN := c.Param("PAN")
	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)
	assetsList, err := contract.SubmitTransaction("GetSupplierVouchers")
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "application/json", assetsList)

	gw.Close()
	client.Close()
}

func setConnection(PAN string) (*grpc.ClientConn, *client.Gateway) { // connect() but for the vouchers
	mspPath := users_common_path + PAN + "/msp"
	certPath := mspPath + "/signcerts/cert.pem"
	keyPath := mspPath + "/keystore"
	peer := "tbchlfdevpeer01"
	domain := "tally.tallysolutions.com"
	peer_port := "7051"
	peerEndpoint := peer + "." + domain + ":" + peer_port
	gatewayPeer := peer + "." + domain
	tlsCertPath := "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/" + domain + "/peers/" + peer + "/tls/ca.crt"

	return connect(peerEndpoint, certPath, keyPath, tlsCertPath, gatewayPeer)
}
