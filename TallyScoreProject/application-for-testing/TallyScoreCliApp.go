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
	cryptoPath   = peer_home + domain 
	TalyScoreccName       = "tallyscore"
	ccName = "tallyscore"
	BusinessProfileCCName="businessprofile"
	channelName  = "tally"

	// CERTPATH, KEYPATH
)

var peer string
var peerEndpoint string
var gatewayPeer string
var tlsCertPath string 

func printUsage()  {
	panic("Usage: \n" +
	"      TallyScoreCliApp <peer_node> register <licenseID>\n" +           
	"      TallyScoreCliApp <peer_node> read <licenseID>\n" +
	"      TallyScoreCliApp <peer_node> increment <licenseID> <inc_by>\n" +
	"      TallyScoreCliApp <peer_node> decrement <licenseID> <dec_by> \n" +
	"      TallyScoreCliApp <peer_node> unregister <licenseID>\n" +
	"      TallyScoreCliApp <peer_node> createVoucher <Voucher_ID> <Supplier_ID> <Voucher_Type> <Hashcode> <TotalValue> <Currency>\n" +
	"      TallyScoreCliApp <peer_node> readVoucher <Voucher_ID>" +
	"      TallyScoreCliApp <peer_node> cancelVoucher <Voucher_ID>" +
	"      TallyScoreCliApp <peer_node> approveVoucher <Voucher_ID>" +
	"      TallyScoreCliApp <peer_node> rejectVoucher <Voucher_ID>" +
	"      TallyScoreCliApp <peer_node> updateVoucher <Voucher_ID> <Value_to_change: Hash or Value> <New_Value>" +
	"      TallyScoreCliApp <peer_node> sendBackVoucher <Voucher_ID>" +
	"      TallyScoreCliApp <peer_node> getSupplierVouchers" +
	"      TallyScoreCliApp <peer_node> getOwnerVouchers" +
	"\n"+
	"  Where:\n" +
	"      <peer_node>: peer host name\n" +
	"      <licenseID> : Company's license ID\n" +
	"      <inc_by>   : increment by how much value\n" +
	"      <dec_by>   : decrement by how much value\n")
}

var user string
var certPath string
var keyPath string

func main(){

	user= os.Getenv("userid") // getenv varaible ---> os.Getenv(userid)
	if user==""{	// in order to run tallyscorecc
		user="Admin"
	}
	certPath= users_common_path + "/" + user + "/msp/signcerts/cert.pem"
	keyPath= users_common_path + "/" + user + "/msp/keystore/"

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
	if ops == "createVoucher" && len(os.Args) < 8 {
		printUsage()
	}
	if ops == "readVoucher" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "cancelVoucher" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "approveVoucher" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "rejectVoucher" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "updateVoucher" && len(os.Args) < 5 {
		printUsage()
	}
	if ops == "sendBackVoucher" && len(os.Args) < 3 {
		printUsage()
	}
	if ops == "getSupplierVouchers" && len(os.Args) < 2 {
		printUsage()
	}
	if ops == "getOwnerVouchers" && len(os.Args) < 2 {
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
	} else if ops == "increment" {
		licenseId := os.Args[3]
		incValue:= os.Args[4]
		fmt.Printf("Incrementing asset of company by %s \n", incValue)
		client, gw:= connect()
		contract := getContract(gw, ccName)
		incrementCompanyScore(contract,licenseId, incValue)
		gw.Close()
		client.Close()
	} else if ops == "decrement" {
		licenseId := os.Args[3]
		decValue:= os.Args[4]
		fmt.Printf("Decrementing asset of company by %s \n", decValue)
		client, gw:= connect()
		contract := getContract(gw, ccName)
		decrementCompanyScore(contract,licenseId, decValue)
		gw.Close()
		client.Close()
	} else if ops == "unregister" {
		licenseId := os.Args[3]
		fmt.Printf("Unregistering company \n")
		client, gw:= connect()
		contract := getContract(gw, ccName)
		unregisterCompany(contract,licenseId)
		gw.Close()
		client.Close()
	} else if ops == "createVoucher" {
		VoucherID:= os.Args[3]
		SupplierID:= os.Args[4]
		VoucherType:= os.Args[5]
		Hashcode:= os.Args[6]
		TotalValue:= os.Args[7]
		Currency:= os.Args[8]
		fmt.Printf("Creating voucher... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		createVoucher(contract, VoucherID, SupplierID, VoucherType, Hashcode, TotalValue, Currency)
		gw.Close()
		client.Close()
	} else if ops == "readVoucher" {
		VoucherID:= os.Args[3]
		fmt.Printf("Reading voucher... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		readVoucher(contract, VoucherID)
		gw.Close()
		client.Close()
	} else if ops == "cancelVoucher" {
		VoucherID:= os.Args[3]
		fmt.Printf("Cancelling voucher... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		cancelVoucher(contract, VoucherID)
		gw.Close()
		client.Close()
	} else if ops == "approveVoucher" {
		VoucherID:= os.Args[3]
		fmt.Printf("Approving voucher... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		approveVoucher(contract, VoucherID)
		gw.Close()
		client.Close()
	} else if ops == "rejectVoucher" {
		VoucherID:= os.Args[3]
		fmt.Printf("Rejecting voucher... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		rejectVoucher(contract, VoucherID)
		gw.Close()
		client.Close()
	} else if ops == "updateVoucher" {
		VoucherID:= os.Args[3]
		toChange:= os.Args[4]
		newValue:= os.Args[5]
		fmt.Printf("Rejecting voucher... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		updateVoucher(contract, VoucherID, toChange, newValue)
		gw.Close()
		client.Close()
	} else if ops == "sendBackVoucher" {
		VoucherID:= os.Args[3]
		fmt.Printf("Sending back voucher... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		sendBackVoucher(contract, VoucherID)
		gw.Close()
		client.Close()
	} else if ops == "getSupplierVouchers" {
		fmt.Printf("Getting Vouchers with you as a supplier... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		getSupplierVouchers(contract)
		gw.Close()
		client.Close()
	} else if ops == "getOwnerVouchers" {
		fmt.Printf("Getting Vouchers with you as an owner... \n")
		client, gw:= connect()
		contract := getContract(gw, BusinessProfileCCName)
		getOwnerVouchers(contract)
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

func unregisterCompany(contract *client.Contract, licenseId string){
	fmt.Printf("\n--> Initiating unregistration of Company: %s\n", licenseId)
	result, err := contract.SubmitTransaction("UnregisterCompany", licenseId)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)
}

func createVoucher(contract *client.Contract, VoucherID string, SupplierID string, VoucherType string, Hashcode string, TotalValue string, Currency string){

	fmt.Printf("\n--> Initiating creation of voucher of user: %s\n", user)
	result, err := contract.SubmitTransaction("VoucherCreated", VoucherID, SupplierID, VoucherType, Hashcode, TotalValue, Currency)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

}

func approveVoucher(contract *client.Contract, VoucherID string){

	fmt.Printf("\n--> Initiating approval of voucher by user: %s\n", user)
	result, err := contract.SubmitTransaction("VoucherApproved", VoucherID)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

}

func rejectVoucher(contract *client.Contract, VoucherID string){

	fmt.Printf("\n--> Initiating rejection of voucher by user: %s\n", user)
	result, err := contract.SubmitTransaction("VoucherRejected", VoucherID)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

}

func sendBackVoucher(contract *client.Contract, VoucherID string){

	fmt.Printf("\n--> Initiating sending back of voucher by user: %s\n", user)
	result, err := contract.SubmitTransaction("VoucherSentBack", VoucherID)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

}

func readVoucher(contract *client.Contract, VoucherID string){
	
	evaluateResult, err := contract.EvaluateTransaction("ReadVoucher", VoucherID) 
	if err != nil {
		fmt.Printf("\n--> Error in reading Voucher's Asset => %s\n", err)
		return
	}
	fmt.Printf("\n--> Voucher details : %s\n", string(evaluateResult))
}

func cancelVoucher(contract *client.Contract, VoucherID string){

	fmt.Printf("\n--> Initiating cancellation of voucher %s.\n", VoucherID)
	result, err := contract.SubmitTransaction("VoucherCancelled", VoucherID)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

}

func updateVoucher(contract *client.Contract, VoucherID string, toChange string, newValue string){

	fmt.Printf("\n--> Initiating updation of voucher %s.\n", VoucherID)
	result, err := contract.SubmitTransaction("VoucherUpdated", VoucherID, toChange, newValue)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)

}

func getSupplierVouchers(contract *client.Contract ) {
	transactionResult, err := contract.EvaluateTransaction("GetVouchersUnderSupplier")
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(transactionResult), err)
}

func getOwnerVouchers(contract *client.Contract ) {
	transactionResult, err := contract.EvaluateTransaction("GetOwnerVouchers")
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(transactionResult), err)
}

func readCompanyAsset(contract *client.Contract, licenseId string){
	evaluateResult, err := contract.EvaluateTransaction("ReadCompanyAsset", licenseId) 
	if err != nil {
		fmt.Printf("\n--> Error in reading Company Asset => %s\n", err)
		return
	}
	fmt.Printf("\n--> Company asset read : %s\n", string(evaluateResult))
}

func incrementCompanyScore(contract *client.Contract, licenseId string, incValue string){
	fmt.Printf("\n--> Initiating increment of score of Company: %s\n", licenseId)
	result, err := contract.SubmitTransaction("IncreaseScore", licenseId, incValue)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)
}

func decrementCompanyScore(contract *client.Contract, licenseId string, decValue string){
	fmt.Printf("\n--> Initiating decrement of score of Company: %s\n", licenseId)
	result, err := contract.SubmitTransaction("DecreaseScore", licenseId, decValue)
	fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)
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

