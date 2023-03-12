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


func printUsage()  {
	panic("Usage: \n" +
	"      integerKeyApp <peer_node> new <var_name>\n" +
	"      integerKeyApp <peer_node> read <var_name>\n" +
	"      integerKeyApp <peer_node> inc <var_name> <inc_by>\n" +
	"      integerKeyApp <peer_node> dec <var_name> <dec_by>\n" +
	"      integerKeyApp <peer_node> del <var_name>\n" +
	"      integerKeyApp <peer_node> list<\n" +
	"\n"+
	"  Where:\n" +
	"      <peer_node>: peer host name\n" +
	"      <var_name> : Variable name\n" +
	"      <inc_by>   : increment by how much value\n" +
	"      <dec_by>   : decrement by how much value\n")
}
func main() {

    if len(os.Args) < 2 {
		printUsage()
    }

	peer = os.Args[1]
	peerEndpoint = peer + "." + domain + ":" + peer_port
	gatewayPeer  = peer + "." + domain
	tlsCertPath  = cryptoPath + "/peers/" + peer + "/tls/ca.crt"

	ops := os.Args[2]

	if ops == "new" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "read" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "del" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "inc" && len(os.Args) < 4 {
		printUsage()
	}
	if ops == "dec" && len(os.Args) < 4 {
		printUsage()
	}

	if ops == "new" {
       var_name := os.Args[3]
	   fmt.Printf("Creating new variable %s \n", var_name)
	   client,gw := connect()
	   contract := getContract(gw)
	   createAsset(contract,var_name)
	   gw.Close()
	   client.Close()
	}else if ops == "read" {
		var_name := os.Args[3]
		fmt.Printf("Reading variable %s \n", var_name)
		client,gw := connect()
		contract := getContract(gw)
		readAsset(contract,var_name)
		gw.Close()
		client.Close()
	  }else if ops == "del" {
		var_name := os.Args[3]
		fmt.Printf("Deleting variable %s \n", var_name)
		client,gw := connect()
		contract := getContract(gw)
		deleteAsset(contract,var_name)
		gw.Close()
		client.Close()
	  }else if ops == "inc" {
		var_name := os.Args[3]
		inc_by := os.Args[4]
		fmt.Printf("Incrementing variable %s by %s\n", var_name, inc_by)
		client,gw := connect()
		contract := getContract(gw)
		increaseValue(contract,var_name, inc_by)
		gw.Close()
		client.Close()
	 }else if ops == "dec" {
		var_name := os.Args[3]
		dec_by := os.Args[4]
		fmt.Printf("Decrementing variable %s by %s\n", var_name, dec_by)
		client,gw := connect()
		contract := getContract(gw)
		decreaseValue(contract,var_name, dec_by)
		gw.Close()
		client.Close()
	}else if ops == "list" {
		fmt.Printf("Listing all variables\n")
		client,gw := connect()
		contract := getContract(gw)
		getAllAssets(contract)
		gw.Close()
		client.Close()
 
	}else{
		printUsage()
	}

}

func connect() (*grpc.ClientConn, *client.Gateway) {
	fmt.Printf("\nConnecting to : %s \n", peerEndpoint)

	// gRPC client conn- shared by all gateway connections to this endpoint
	clientConnection := newGrpcConnection()

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
	
	return clientConnection, gw


}

func getContract(gw *client.Gateway ) *client.Contract {
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

// function to call the ReadAsset function present in smartcontract.go
func readAsset(contract *client.Contract , name string) {

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", name) // EvaluateTransaction evaluates a transaction in the scope of the specified context and returns its context
	if err != nil {
		fmt.Printf("\n--> Error in reading Asset : %s => %s\n", name, err)
		return
	}
	fmt.Printf("\n--> Read Asset : %s => %s\n", name, string(evaluateResult))

}

func createAsset(contract *client.Contract , name string) {

	fmt.Printf("\n--> Creating Asset : %s\n", name)

	result, err := contract.SubmitTransaction("CreateAsset", name) // SubmitTransaction returns results of a transaction only after its commited

	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

}

func increaseValue(contract *client.Contract , name string, incVal string) {

	fmt.Printf("Name : %s , IncreaseValue: %s ", name, incVal)

	evaluatedAsset, err := contract.SubmitTransaction("IncreaseAsset", name, incVal)
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(evaluatedAsset), err)
}

func decreaseValue(contract *client.Contract , name string, decVal string) {

	fmt.Printf("Name : %s , DecreaseValue: %s ", name, decVal)

	evaluatedAsset, err := contract.SubmitTransaction("DecreaseAsset", name, decVal)
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(evaluatedAsset), err)
}


func getAllAssets(contract *client.Contract ) {

	transactionResult, err := contract.EvaluateTransaction("GetAllAssets")

	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(transactionResult), err)

}

func deleteAsset(contract *client.Contract , name string){


	_, err := contract.SubmitTransaction("DeleteAsset", name) 
	fmt.Printf("\n------> After SubmitTransaction: %s \n",  err)
}


func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return "error in parsing JSON"
	}
	return prettyJSON.String()
}
