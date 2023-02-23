package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "/home/hlfabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

type CreateAssetRequest struct {                         // for create operations
	Name string `json:"name" binding:"required"`
}


type UpdateValueRequest struct {
	Name         string `json:"Name" binding:"required"`
	Value string `json:"Value" binding:"required"`
}

var contract *client.Contract

func main() {
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
	chaincodeName := "integerKey"
	channelName := "tallychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract = network.GetContract(chaincodeName)

	router := gin.Default()

	router.PUT("/integerKey/createAsset", createAsset)
	router.GET("/integerKey/readAsset/:name", readAsset)
	router.POST("/integerKey/increaseValue", increaseValue)
	router.POST("/integerKey/decreaseValue", decreaseValue)
	router.GET("/integerKey/getAllAssets", getAllAssets)
	router.DELETE("/integerKey/deleteAsset/:name", deleteAsset)
	router.GET("integerKey/getPagination/:startName/:endName", getPagination)
	router.Run("localhost:8080")

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

		c.String(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}\n", err))

	}

	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluateResult)))
}

func createAsset(c *gin.Context) {

	var request CreateAssetRequest
	c.BindJSON(&request)
	name := request.Name

	fmt.Printf("\n--> Creating Asset : %s\n", name)

	result, err := contract.SubmitTransaction("CreateAsset", name) // SubmitTransaction returns results of a transaction only after its commited

	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}\n", err))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("{\"name\":\"%s\",\"value\":\"0\"}\n", name))

}

func increaseValue(c *gin.Context) {

	var request UpdateValueRequest
	c.BindJSON(&request)
	name := request.Name
	incVal := request.Value

	fmt.Printf("Name : %s , IncreaseValue: %s ", name, incVal)

	evaluatedAsset, err := contract.SubmitTransaction("IncreaseAsset", name, incVal)
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(evaluatedAsset), err)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}\n", err))
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluatedAsset)))

}

func decreaseValue(c *gin.Context) {

	var request UpdateValueRequest
	c.BindJSON(&request)
	name := request.Name
	decVal := request.Value

	fmt.Printf("Name : %s , DecreaseValue: %s ", name, decVal)

	evaluatedAsset, err := contract.SubmitTransaction("DecreaseAsset", name, decVal)
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(evaluatedAsset), err)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}\n", err))
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluatedAsset)))
}


func getAllAssets(c *gin.Context) {

	transactionResult, err := contract.EvaluateTransaction("GetAllAssets")

	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("%s\n", string(transactionResult)))

}


func getPagination(c *gin.Context){

	startname := c.Param("startname")
	endname := c.Param("endname")
	// pageSize := c.Param("pageSize")
	transactionResult, err := contract.EvaluateTransaction("GetAssetsPagination", startName, endName, int32(5), string(""))
	if err != nil{
		return
	}

	c.IndentedJSON(http.StatusOK, fmt.Sprintf("%s \n", string(transactionResult)))

}






func deleteAsset(c *gin.Context){

	name := c.Param("name")

	_, err := contract.SubmitTransaction("DeleteAsset", name) 
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}\n", err))
	}
	//c.String(http.StatusOK, fmt.Sprintf("%s\n", string(deleteAssetResult)))
	//c.String(http.StatusOK, deleteAssetResult)
	c.JSON(http.StatusOK, gin.H{name:"has been deleted"})
}


func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return "error in parsing JSON"
	}
	return prettyJSON.String()
}