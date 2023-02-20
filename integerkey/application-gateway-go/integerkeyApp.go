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

	// createAsset(contract, "key1")
	// increaseValue(contract, "key1", 5) // we want to increase the value of the asset by 5
	// decreaseValue(contract, "key1", 2)
	// readAsset(contract, "key1")
	// readAsset(contract, "key2")
	// fmt.Print("TESTING DONE")

	router := gin.Default()

	router.GET("/integerKey/createAsset/:name", createAsset)
	router.GET("/integerKey/readAsset/:name", readAsset)

	router.POST("/integerKey/increaseValue/:name/:value", increaseValue)
	router.POST("/integerKey/decreaseValue/:name/:value", decreaseValue)

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

func createAsset(c *gin.Context) {

	name := c.Param("name")
	evaluateResult, err := contract.EvaluateTransaction("CreateAsset", name)
	if err != nil {
		// panic(fmt.Errorf("failed to evaluate transaction: %w", err))

		c.IndentedJSON(http.StatusNotImplemented, gin.H{"message": "failed to evaluate transaction"})

	}
	result := formatJSON(evaluateResult)

	c.IndentedJSON(http.StatusOK, result)

}

func increaseValue(c *gin.Context) {
	// fmt.Printf("\n--> Submit Transaction: Increase Asset Value (by %v) \n", incVal)
	name := c.Param("name")
	incVal := c.Param("value")

	evaulateResult, err := contract.SubmitTransaction("IncreaseAsset", name, incVal)
	if err != nil {
		// panic(fmt.Errorf("failed to submit transaction: %w", err))
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"message": "failed to evaluate transaction"})
	}

	result := formatJSON(evaulateResult)

	c.IndentedJSON(http.StatusOK, result)
	//fmt.Printf("*** Transaction committed successfully\n")
}
func decreaseValue(c *gin.Context) {
	// fmt.Printf("\n--> Submit Transaction: Decrease Asset Value (by %v) \n", incVal)
	name := c.Param("name")
	decVal := c.Param("value")

	// evaulateResult, err := contract.SubmitTransaction("DecreaseAsset", name, strconv.FormatUint(uint64(decVal), 10))
	evaulateResult, err := contract.SubmitTransaction("DecreaseAsset", name, decVal)
	if err != nil {
		// panic(fmt.Errorf("failed to submit transaction: %w", err))
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"message": "failed to evaluate transaction"})
	}

	result := formatJSON(evaulateResult)

	c.IndentedJSON(http.StatusOK, result)
	//fmt.Printf("*** Transaction committed successfully\n")
}

func readAsset(c *gin.Context) {
	// fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	name := c.Param("name")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", name)
	if err != nil {
		// panic(fmt.Errorf("failed to evaluate transaction: %w", err))

		c.IndentedJSON(http.StatusNotImplemented, gin.H{"message": "failed to evaluate transaction"})

	}
	result := formatJSON(evaluateResult)

	c.IndentedJSON(http.StatusOK, result)

	// fmt.Printf("*** Result:%s\n", result)
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return "error in parsing JSON"

	}
	return prettyJSON.String()
}
