package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/itsjamie/gin-cors"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"net/http"

	"github.com/gin-gonic/gin"
)



const (
	mspID        = "Tally"
	peer_home    = "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/"
	domain       = "tally.tallysolutions.com"
	user         = "Admin"
	peer_port    = "7051"
	cryptoPath   = peer_home + domain 
	certPath     = cryptoPath + "/users/" + user +  "@" + domain + "/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/" + user +  "@" + domain + "/msp/keystore/"
	ccName       = "integerkey"
	channelName  = "integerkey"
)

var peer string
var peerEndpoint string
var gatewayPeer string
var tlsCertPath string 

type CreateAssetRequest struct {                         // for create operations
	Name string `json:"name" binding:"required"`
	Owner string `json:"owner" binding:"required"`
}


type UpdateValueRequest struct {
	Name  string `json:"Name" binding:"required"`
	Value string `json:"Value" binding:"required"`
	Owner string `json:"Owner" binding:"required"`
}

var contract *client.Contract

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	peer = hostname
	peerEndpoint = peer + "." + domain + ":" + peer_port
	gatewayPeer  = peer + "." + domain
	tlsCertPath  = cryptoPath + "/peers/" + peer + "/tls/ca.crt"


	// gRPC client conn- shared by all gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	//creating client identity, signing implementation
	id := newIdentity()
	sign := newSign()

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
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract = network.GetContract(ccName)

	router := gin.New()

	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge: 50 * time.Second,
		Credentials: false,
		ValidateHeaders: false,
	}))

	router.PUT("/integerKey/createAsset", createAsset)
	router.GET("/integerKey/readAsset/:name", readAsset)
	router.POST("/integerKey/increaseValue", increaseValue)
	router.POST("/integerKey/decreaseValue", decreaseValue)
	router.GET("/integerKey/getAllAssets", getAllAssets)
	router.DELETE("/integerKey/deleteAsset/:name", deleteAsset)
	router.GET("/integerKey/getPagination/:startName/:endName/:bookmark", getPagination)
	router.Run("0.0.0.0:8080")

}

func newGrpcConnection() *grpc.ClientConn {
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

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
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

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
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

// function to call the ReadAsset function present in smartcontract.go
func readAsset(c *gin.Context) {

	name := c.Param("name")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", name) // EvaluateTransaction evaluates a transaction in the scope of the specified context and returns its context
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
	}
	c.Writer.Header().Set("Content-Type","application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluateResult)))
}

func createAsset(c *gin.Context) {

	var request CreateAssetRequest
	c.BindJSON(&request)
	name := request.Name
	owner:= request.Owner

	fmt.Printf("\n--> Creating Asset : %s with owner: %s\n", name, owner)

	result, err := contract.SubmitTransaction("CreateAsset", name, owner) // SubmitTransaction returns results of a transaction only after its commited

	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin","*")
	c.Writer.Header().Set("Access-Control-Allow-Methods","PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers","Content-Type")
	// c.String(http.StatusOK, fmt.Sprintf("{\"name\":\"%s\",\"value\":\"0\"}\n", name))

	c.JSON(http.StatusOK, gin.H{"Name":name, "Value":0, "Owner":owner})

}

func increaseValue(c *gin.Context) {

	var request UpdateValueRequest
	c.BindJSON(&request)
	name := request.Name
	incVal := request.Value
	owner := request.Owner

	fmt.Printf("Name : %s , IncreaseValue: %s , Owner: %s ", name, incVal, owner)

	evaluatedAsset, err := contract.SubmitTransaction("IncreaseAsset", name, incVal, owner)
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(evaluatedAsset), err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	c.Writer.Header().Set("Content-Type","application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluatedAsset)))

}

func decreaseValue(c *gin.Context) {

	var request UpdateValueRequest
	c.BindJSON(&request)
	name := request.Name
	decVal := request.Value
	owner := request.Owner

	fmt.Printf("Name : %s , DecreaseValue: %s , Owner: %s ", name, decVal, owner)

	evaluatedAsset, err := contract.SubmitTransaction("DecreaseAsset", name, decVal, owner)
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(evaluatedAsset), err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	c.Writer.Header().Set("Content-Type","application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluatedAsset)))
}


func getAllAssets(c *gin.Context) {

	transactionResult, err := contract.EvaluateTransaction("GetAllAssets")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Content-Type","application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(transactionResult)))
	//c.JSON(http.StatusOK, gin.H{})

}


// gets the assets between a specified range 
func getPagination(c *gin.Context){

	startname := c.Param("startname")
	endname := c.Param("endname")
	bookmark := c.Param("bookmark")
	// pageSize := c.Param("pageSize") -- PAGE SIZE IS NOT PASSED AS A PARAMETER (FOR NOW- will consider if required in later use cases)
	transactionResult, err := contract.EvaluateTransaction("GetAssetsPagination", startname, endname, bookmark)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	} 
	c.Writer.Header().Set("Content-Type","application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(transactionResult)))

}


func deleteAsset(c *gin.Context){

	name := c.Param("name")

	_, err := contract.SubmitTransaction("DeleteAsset", name) 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	//c.String(http.StatusOK, fmt.Sprintf("%s\n", string(deleteAssetResult)))
	//c.String(http.StatusOK, deleteAssetResult)
	c.JSON(http.StatusOK, gin.H{name:"has been deleted"})
}


// func clearAllAssets(c *gin.Context){
// 	//[]* allAssets =
// }


func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return "error in parsing JSON"
	}
	return prettyJSON.String()
}


