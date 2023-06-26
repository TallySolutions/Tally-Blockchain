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
	"github.com/itsjamie/gin-cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	mspID       = "Tally"
	peer_home   = "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/"
	domain      = "tally.tallysolutions.com"
	user        = "Admin"
	peer_port   = "7051"
	cryptoPath  = peer_home + domain
	certPath    = cryptoPath + "/users/" + user + "@" + domain + "/msp/signcerts/cert.pem"
	keyPath     = cryptoPath + "/users/" + user + "@" + domain + "/msp/keystore/"
	ccName      = "integerkey"
	channelName = "integerkey"
)

var peer string
var peerEndpoint string
var gatewayPeer string
var tlsCertPath string

type CreateOrganizationRequest struct {
	OrgId string `json:"org_id" binding:"required"`
}

type CreateElectionRequest struct {
	ElectionId string `json:"election_id" binding:"required"`
}

type CreateVoteOptionsRequest struct {
	Options []string `json:"options" binding:"required"`
}

type CastVoteRequest struct {
	UserId string   `json:"user_id" binding:"required"`
	Picks  []string `json:"picks" binding:"required"`
}

type VoteOptionsResponse struct {
	Options []string `json:"options" binding:"required"`
}

type Member struct {
	MemberId string `json:"member_id" binding:"required"`
	Role     string `json:"role" binding:"required"`
}
type CreateMembersRequest struct {
	Memebers []Member `json:"members" binding:"required"`
}

type VoteResult struct {
	OptionId string `json:"option_id" binding:"required"`
	Count    int    `json:"count" binding:"required"`
}
type VoteResultResponse struct {
	Options []VoteResult `json:"options" binding:"required"`
}

type Ballot struct {
	CastedBy  string `json:"casted_by" binding:"required"`
	Timestamp string `json:"casted_at" binding:"required"`
}

type DetailedVoteResult struct {
	OptionId string   `json:"option_id" binding:"required"`
	Count    int      `json:"count" binding:"required"`
	Ballots  []Ballot `json:"ballots" binding:"required"`
}

type DetailedVoteResultResponse struct {
	Options []DetailedVoteResult `json:"options" binding:"required"`
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
	gatewayPeer = peer + "." + domain
	tlsCertPath = cryptoPath + "/peers/" + peer + "/tls/ca.crt"

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

	// handle owner contract

	router := gin.New()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, POST",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	//Orgqanizations
	router.POST("/smartevm/organization/create", createOrganization)

	//Election
	router.POST("/smartevm/:organization/election/create", createElection)

	//Vote
	router.POST("/smartevm/:organization/:election/options/add", addVotingOptions)
	router.POST("/smartevm/:organization/:election/members/add", addMembers)
	router.POST("/smartevm/:organization/:election/cast", castVote)
	router.GET("/smartevm/:organization/:election/options", getAllOptions)
	router.GET("/smartevm/:organization/:election/results", getResults)
	router.GET("/smartevm/:organization/:election/results/detailed", getDetailedResults)

	//Start router
	router.Run("0.0.0.0:8080")

}

func respondWithError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
}

// Organization
func createOrganization(c *gin.Context) {

	var request CreateOrganizationRequest
	c.BindJSON(&request)
	orgId := request.OrgId

	fmt.Printf("\n--> Creating Organization : %s\n", orgId)

	//TODO: Create HL Fabric Organization here ...

	respondWithError(c, fmt.Errorf("Method not implemented."))
}

// Election
func createElection(c *gin.Context) {

	org_id := c.Param("organization")
	var request CreateElectionRequest
	c.BindJSON(&request)
	electionId := request.ElectionId

	fmt.Printf("\n--> Creating Channel '%s' for organization %s\n", electionId, org_id)

	//TODO: Create HL Fabric Channel for Organization here ...

	respondWithError(c, fmt.Errorf("Method not implemented."))
}

// Vote
func addVotingOptions(c *gin.Context) {
	org_id := c.Param("organization")
	election_id := c.Param("election")
	fmt.Printf("\n--> Adding voting options for %s -> %s\n", org_id, &election_id)
	var request CreateVoteOptionsRequest
	c.BindJSON(&request)

	//TODO: Implement add voting ..

	for i, option := range request.Options {
		fmt.Printf("   %d) %s\n", i, option)
	}

	respondWithError(c, fmt.Errorf("Method not implemented."))

}
func addMembers(c *gin.Context) {
	org_id := c.Param("organization")
	election_id := c.Param("election")
	fmt.Printf("\n--> Adding members for %s -> %s\n", org_id, &election_id)
	var request CreateMembersRequest
	c.BindJSON(&request)

	//TODO: Implement add members ..

	for i, member := range request.Memebers {
		fmt.Printf("   %d) %s (%s)\n", i, member.MemberId, member.Role)
	}

	respondWithError(c, fmt.Errorf("Method not implemented."))

}
func castVote(c *gin.Context) {
	org_id := c.Param("organization")
	election_id := c.Param("election")
	fmt.Printf("\n--> Casting vote for %s -> %s\n", org_id, &election_id)
	var request CastVoteRequest
	c.BindJSON(&request)

	//TODO: Implement cast vote ...

	fmt.Printf("   %s :\n", request.UserId)

	for i, pick := range request.Picks {
		fmt.Printf("      %d) %s\n", i, pick)
	}

	respondWithError(c, fmt.Errorf("Method not implemented."))
}
func getAllOptions(c *gin.Context) {
	org_id := c.Param("organization")
	election_id := c.Param("election")
	fmt.Printf("\n--> Listing all options for %s -> %s\n", org_id, &election_id)

	//TODO: Implement get all options ...

	respondWithError(c, fmt.Errorf("Method not implemented."))
}
func getResults(c *gin.Context) {
	org_id := c.Param("organization")
	election_id := c.Param("election")
	fmt.Printf("\n--> Getting results for %s -> %s\n", org_id, &election_id)

	//TODO: Implement get results ..

	respondWithError(c, fmt.Errorf("Method not implemented."))

}
func getDetailedResults(c *gin.Context) {
	org_id := c.Param("organization")
	election_id := c.Param("election")
	fmt.Printf("\n--> Getting detailed results for %s -> %s\n", org_id, &election_id)
	//TODO: Implement get detaailed results ..

	respondWithError(c, fmt.Errorf("Method not implemented."))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, fmt.Sprintf("%s\n", string(evaluateResult)))
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return "error in parsing JSON"
	}
	return prettyJSON.String()
}
