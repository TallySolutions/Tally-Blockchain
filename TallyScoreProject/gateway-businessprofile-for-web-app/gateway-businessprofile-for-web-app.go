package main

// todo before running this code: start the CA servers

// contains gateway code for defining API endpoints for voucher related stuff

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path"
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
	BusinessProfileCCName = "businessprofile"
	channelname           = "tally"
	users_common_path     = "/home/ubuntu/fabric/tally-network/fabric-ca-servers/tally/client/users/"
)

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
		timeoutGroup.PUT("/TallyScoreProject/voucherCreation/:PAN", voucherCreation)
		timeoutGroup.POST("/TallyScoreProject/voucherCancellation/:PAN", voucherCancellation)
		timeoutGroup.PUT("/TallyScoreProject/voucherUpdation/:PAN", voucherUpdation)
		timeoutGroup.GET("/TallyScoreProject/listOwnerVouchers/:PAN", listOwnerVouchers)

		timeoutGroup.POST("/TallyScoreProject/voucherApproval/:PAN", voucherApproval)
		timeoutGroup.POST("/TallyScoreProject/voucherReturn/:PAN", voucherReturn)
		timeoutGroup.POST("/TallyScoreProject/voucherRejection/:PAN", voucherRejection)
		timeoutGroup.GET("/TallyScoreProject/listSupplierVouchers/:PAN", listSupplierVouchers)
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

	c.JSON(http.StatusOK, gin.H{"message": "Voucher created successfully"})
}
func voucherCancellation(c *gin.Context) {

	var request voucherIDRequest
	c.BindJSON(&request)
	VoucherID := request.VoucherID
	PAN := c.Param("PAN")

	fmt.Printf("Initiating cancellation of Voucher %s", VoucherID)
	//(Display voucher info)(Use ReadVoucher function)
	client, gw := setConnection(PAN)
	network := gw.GetNetwork(channelname)
	contract := network.GetContract(BusinessProfileCCName)

	result, err := contract.SubmitTransaction("VoucherCancelled", VoucherID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		fmt.Printf("\n Submit Transaction returned: O/p= %s , Error= %s \n", string(result), err)
		return
	}

	asset, err := contract.SubmitTransaction("ReadVoucher", VoucherID)
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "application/json", asset)
	gw.Close()
	client.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Voucher cancelled successfully"})

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
	c.Data(http.StatusOK, "application/json", asset)

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

func setConnection(PAN string) (*grpc.ClientConn, *client.Gateway) {
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