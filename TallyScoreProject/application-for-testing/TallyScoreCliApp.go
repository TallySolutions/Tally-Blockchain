package main

import (
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

)

const (
	mspID        = "Tally"
	peer_home    = "/home/ubuntu/fabric/tally-network/organizations/peerOrganizations/"
	users_common_path = "/home/ubuntu/fabric/tally-network/clients/users"
	domain       = "tally.tallysolutions.com"
	peer_port    = "7051"
	user= "Admin"
	cryptoPath   = peer_home + domain 
	certPath     = cryptoPath + "/users/" + user +  "@" + domain + "/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/" + user +  "@" + domain + "/msp/keystore/"
	TalyScoreccName       = "tallyscore"
	ccName = "tallyscore"
	channelName  = "tally"

	// CERTPATH, KEYPATH
)

var peer string
var peerEndpoint string
var gatewayPeer string
var tlsCertPath string 

func printUsage()  {
	panic("Usage: \n" +
	"      integerKeyApp <peer_node> register <licenseID>\n" +           
	"      integerKeyApp <peer_node> read <licenseID>\n" +
	"      integerKeyApp <peer_node> increment <licenseID> <inc_by>\n" +
	"      integerKeyApp <peer_node> decrement <licenseID> <dec_by> \n" +
	"      integerKeyApp <peer_node> unregister <licenseID>\n" +
	"\n"+
	"  Where:\n" +
	"      <peer_node>: peer host name\n" +
	"      <licenseID> : Company's license ID\n" +
	"      <inc_by>   : increment by how much value\n" +
	"      <dec_by>   : decrement by how much value\n")
}


func main(){

	fmt.Printf("USER:%s \n", user)

    if len(os.Args) < 2 {
		printUsage()
    }

	peer = os.Args[1]
	peerEndpoint = peer + "." + domain + ":" + peer_port
	gatewayPeer  = peer + "." + domain
	tlsCertPath  = cryptoPath + "/peers/" + peer + "/tls/ca.crt"

	ops := os.Args[2]
	fmt.Printf("ops: %s\n", ops)

	if ops == "register" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "read" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "increment" && len(os.Args) < 4 {
		printUsage()
	}
	if ops == "decrement" && len(os.Args) < 4 {
		printUsage()
	}
	if ops == "unregister" && len(os.Args) < 3 {
		printUsage()
	}
	
	if ops == "register" {
		licenseId := os.Args[3]
		fmt.Printf("Initiating registration of the company %s \n", licenseId)
		client, gw := connect()
		contract := getContract(gw, ccName)
		registerCompany(contract, licenseId)
		gw.Close()
		client.Close()
	 } else if ops == "read" {
		licenseId := os.Args[3]
		fmt.Printf("Reading asset of company with id: %s \n", licenseId)
		client, gw:= connect()
		contract := getContract(gw, ccName)
		readCompanyAsset(contract,licenseId)
		gw.Close()
		client.Close()
	} else{
		printUsage()
	}
}

func registerCompany(contract *client.Contract, licenseId string){
	fmt.Printf("\n--> Initiating registration of Company: %s\n", licenseId)
	result, err := contract.SubmitTransaction("RegisterCompany", licenseId)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)
}

func readCompanyAsset(contract *client.Contract, licenseId string){
	evaluateResult, err := contract.EvaluateTransaction("ReadCompanyAsset", licenseId) 
	if err != nil {
		fmt.Printf("\n--> Error in reading Company Asset => %s\n", err)
		return
	}
	fmt.Printf("\n--> Company asset read : %s\n", string(evaluateResult))
}

func connect() (*grpc.ClientConn, *client.Gateway) {
	fmt.Printf("\nConnecting to : %s \n", peerEndpoint)

	// gRPC client conn- shared by all gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	//creating client identity, signing implementation
	id := newIdentity()              // stores client id
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
	return clientConnection, gw
}

func getContract(gw *client.Gateway , ccName string ) *client.Contract {
	network := gw.GetNetwork(channelName)
	return  network.GetContract(ccName)
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

