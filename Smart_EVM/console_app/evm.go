package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)


const (
	mspID        = "VoterMSP"
	peer_home    = "/home/ubuntu/fabric/evm-network/organizations/peerOrganizations/"
	domain       = "voter.boatinthesea.com"
	user         = "Admin"
	peer_port    = "9051"
	cryptoPath   = peer_home + domain
	certPath     = cryptoPath + "/users/" + user +  "@" + domain + "/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/" + user +  "@" + domain + "/msp/keystore/"
	ccName       = "evm"
	channelName  = "tally"

)


var (
	peer         = "peer0"
	peerEndpoint = "localhost:" + peer_port
	gatewayPeer  = peer + "." + domain
	tlsCertPath  = cryptoPath + "/peers/" + peer + "." + domain + "/tls"
)


func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <Command> [OPTIONS]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "\tinit : Initialize the ledger, use flags -anon, -abs, -single to set the ledger properties.")
	fmt.Fprintln(os.Stderr, "\toptions : Voting options command, use flag -o to set the options")
	fmt.Fprintln(os.Stderr, "\tvoters : Voters command, use flag -v to set the voters")
	fmt.Fprintln(os.Stderr, "\tcast : Cast vote, usage: cast <voter> -o=<comma-separated-options>")
	fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
}

// Define a custom flag type that stores multiple values
type stringSliceFlag []string

// Implement the String method for the custom flag type
func (f *stringSliceFlag) String() string {
	return fmt.Sprintf("%v", *f)
}

// Implement the Set method for the custom flag type
func (f *stringSliceFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {

	// Define flags for each option

	//Init flags
	isAnonymous := flag.Bool("anon", false, "Anonymous Voting")
	isAbstainable := flag.Bool("abs", false, "Abstainable Voting")
	isSingle := flag.Bool("single", false, "Single Choice Voting")

	var options stringSliceFlag
	flag.Var(&options, "o", "List of options (comma-separated)")

	var voters stringSliceFlag
	flag.Var(&voters, "v", "List of voters (comma-separated)")

	// Set the custom usage function
	flag.Usage = printUsage

	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	peer = os.Args[1]



	if os.Args[1] == "init" {
		initCommand(*isAnonymous, *isAbstainable, *isSingle)
	} else if os.Args[1] == "options" {
		optionsCommand(options)
	} else if os.Args[1] == "voters" {
		votersCommand(voters)
	} else if os.Args[1] == "cast" {
		if len(os.Args) < 3 {
			flag.Usage()
			os.Exit(1)
		}
		castCommand(os.Args[2], options)
	} else {
		flag.Usage()
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


// define function for each command
func initCommand(isAnonymous bool, isAbstainable bool, isSingle bool) {
	fmt.Println("Initializing the ledger")
	//Call the init ledger blockchain function with the parameters

	//1. Connect to the blockchain
	clientConnection ,gw := connect()
	defer clientConnection.Close()

	//2. Init function
	contract :=getContract(gw)

	//Submit the transaction to initialize the ledger
	_, err := contract.SubmitTransaction("Initledger", strconv.FormatBool(isAnonymous) , strconv.FormatBool(isAbstainable) , strconv.FormatBool(isSingle))
	if err != nil {
		fmt.Println("Error initializing ledger:", err)
	} else {
		fmt.Println("Ledger initialization successful")
	}

}



func optionsCommand(options []string) {
	fmt.Println("Setting options")

	// Connect to the blockchain
	clientConnection ,gw := connect()
	defer clientConnection.Close()

	//set options function with the parameters
	contract:=getContract(gw)

	_, err:= contract.SubmitTransaction("SetOptions" , options...)
	if err!=nil{
		fmt.Println("Error setting options:",err)
	} else{
		fmt.Println("Options set successfully")
	}
}

func votersCommand(voters []string) {
	fmt.Println("Setting voters")

	//1. Connect to the blockchain
	clientConnection ,gw := connect()
	defer clientConnection.Close()

	//set voters function with the parameters
	contract:=getContract(gw)

	_,err :=contract.SubmitTransaction("SetVoters" , voters...)
	if err!=nil{
		fmt.Println("Error setting voters:",err)
	} else{
		fmt.Println("Voters set successfully")
	}

}

func castCommand(voter string, options []string) {
	fmt.Println("Casting vote")

  //1. Connect to the blockchain
	clientConnection ,gw := connect()
	defer clientConnection.Close()

	//cast vote function with the parameters
	contract:= getContract(gw)

	args := append([]string{voter}, options...)

	_,err:=contract.SubmitTransaction("CastVote", args...)
	if err!=nil{
		fmt.Println("Error casting vote:",err)
	} else{
		fmt.Println("Vote cast successfully")
	}
}
