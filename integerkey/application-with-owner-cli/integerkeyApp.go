package main


// HANDLES BOTH OWNER ASSET OPERATIONS, AS WELL AS INTEGER KEY ASSET OPERATIONS



import (

	//"encoding/asn1"
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
	users_common_path = "/home/ubuntu/fabric/tally-network/clients/users"
	domain       = "tally.tallysolutions.com"
	user         = "user2"
	peer_port    = "7051"
	cryptoPath   = peer_home + domain 
	certPath     = users_common_path + "/" + user + "/msp/signcerts/cert.pem"
	keyPath      = users_common_path + "/" + user + "/msp/keystore/"
	intkeyccName       = "integerkey"
	ccName = "integerkey"
	ownerccName = "owner"
	channelName  = "tally"

)

var peer string
var peerEndpoint string
var gatewayPeer string
var tlsCertPath string 


func printUsage()  {
	panic("Usage: \n" +
	"      integerKeyApp <peer_node> reg_owner <owner_name>\n" +
	"      integerKeyApp <peer_node> all_owners\n" +
	"      integerKeyApp <peer_node> unreg_owner <owner_name>\n" +
	"      integerKeyApp <peer_node> del_owner <owner_name>\n" +
	"      integerKeyApp <peer_node> new_asset <var_name>\n" +           
	"      integerKeyApp <peer_node> read <var_name>\n" +
	"      integerKeyApp <peer_node> inc <var_name> <inc_by>\n" +
	"      integerKeyApp <peer_node> transfer_asset <var_name>\n" +
	"      integerKeyApp <peer_node> dec <var_name> <dec_by> \n" +
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
	fmt.Printf("ops: %s\n", ops)

	if ops == "new_asset" && len(os.Args) < 3 {
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

	if ops == "reg_owner" {
       // owner_id := os.Args[3]
	   owner_name := os.Args[3]
	   fmt.Printf("Registering owner %s \n", owner_name)
	   client, gw := connect()
	   ownercontract := getContract(gw, ownerccName)
	   RegisterOwner(ownercontract , owner_name)
	   gw.Close()
	   client.Close()
	} else if ops == "all_owners" {
		fmt.Printf("getting all owners\n")
		client,gw := connect()
		ownercontract := getContract(gw, ownerccName)
		GetAllOwners(ownercontract)
		gw.Close()
		client.Close()
	 } else if ops == "unreg_owner" {
		owner_name := os.Args[3]
		fmt.Printf("Unregistering owner %s \n", owner_name)
		client, gw := connect()
		ownercontract := getContract(gw, ownerccName)
		UnregisterOwner(ownercontract,owner_name)
		gw.Close()
		client.Close()
	 } else if ops == "del_owner" {
		owner_name := os.Args[3]
		fmt.Printf("Deleting owner %s \n", owner_name)
		client, gw := connect()
		contract := getContract(gw, ownerccName)
		deleteOwner(contract , owner_name)
		gw.Close()
		client.Close()
	 } else if ops == "new_asset" {
		var_name := os.Args[3]
		fmt.Printf("Initiating creation of new asset %s \n", var_name)
		client, gw := connect()
		contract := getContract(gw, ccName)
		createAsset(contract, var_name)
		gw.Close()
		client.Close()
	 } else if ops == "read" {
		var_name := os.Args[3]
		fmt.Printf("Reading variable %s \n", var_name)
		client, gw:= connect()
		contract := getContract(gw, ccName)
		readAsset(contract,var_name)
		gw.Close()
		client.Close()
	  }else if ops == "del" {
		var_name := os.Args[3]
		fmt.Printf("Deleting variable %s \n", var_name)
		client, gw:= connect()
		contract := getContract(gw, ccName)
		deleteAsset(contract,var_name)
		gw.Close()
		client.Close()
	  }else if ops == "inc" {
		var_name := os.Args[3]
		inc_by := os.Args[4]
		fmt.Printf("Incrementing variable %s by %s\n", var_name, inc_by)
		client, gw:= connect()
		contract := getContract(gw, ccName)
		increaseValue(contract, var_name, inc_by)
		gw.Close()
		client.Close()
	 }else if ops == "dec" {
		var_name := os.Args[3]
		dec_by := os.Args[4]
		fmt.Printf("Decrementing variable %s by %s\n", var_name, dec_by)
		client, gw := connect()
		contract := getContract(gw, ccName)
		decreaseValue(contract, var_name, dec_by)
		gw.Close()
		client.Close()
	}else if ops == "transfer_asset" {
		var_name := os.Args[3]
		fmt.Printf("Transferring asset %s to user %s \n", var_name, user)
		client, gw := connect()
		contract := getContract(gw, ccName)
		transferAsset(contract, var_name)
		gw.Close()
		client.Close()
	}else if ops == "list" {
		fmt.Printf("Listing all variables\n")
		client, gw := connect()
		contract := getContract(gw, ccName)
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

	// clientID:=id.Mspid
	// idBytes := id.Identifier().Id
	// idString := string(idBytes)
	// fmt.Println("ID:", idString)

	// idBytes := id.Identifier.Id
	// idString := string(idBytes)
	// fmt.Println("ID:", idString)

	fmt.Printf("TYPE OF ID: %T \n", id)
	
	// return clientConnection, gw, string(clientID)
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

	// enrollmentID := ""
    // for _, ext := range certificate.Extensions {
    //     if ext.Id.Equal(asn1.ObjectIdentifier{1,2,3,4,5,6,7,8,1}) { // OID for hf.EnrollmentID
    //         var value string
    //         _, err := asn1.Unmarshal(ext.Value, &value)
    //         if err != nil {
    //             return nil, err
    //         }
    //         enrollmentID = value
    //         break
    //     }
    // }

	// id, err := identity.NewX509Identity(mspID, certificate)   
	// id.EnrollmentID= enrollmentID

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





func RegisterOwner(contract *client.Contract, owner_name string ) {
	evaluateResult, err := contract.SubmitTransaction("RegisterOwner", owner_name) // EvaluateTransaction evaluates a transaction in the scope of the specified context and returns its context
	if err != nil {
		fmt.Printf("\n--> Error in reading Asset : %s => %s\n", owner_name, err)
		return
	}
	fmt.Printf("\n--> Registered owner : %s . %s\n", owner_name, string(evaluateResult))
}

func UnregisterOwner(contract *client.Contract, owner_name string ) {
	evaluateResult, err := contract.SubmitTransaction("UnregisterOwner", owner_name) // EvaluateTransaction evaluates a transaction in the scope of the specified context and returns its context
	if err != nil {
		fmt.Printf("\n--> Error in reading Asset : %s => %s\n", owner_name, err)
		return
	}
	fmt.Printf("\n--> Unregistered owner : %s . %s\n",owner_name, string(evaluateResult))
}
func GetAllOwners(contract *client.Contract ) {

	transactionResult, err := contract.EvaluateTransaction("GetAllOwners")

	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(transactionResult), err)

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

// func createAsset(contract *client.Contract, name string) {

// 	owner_exists, err := ownercontract.EvaluateTransaction("OwnerExistence", ownername)
// 	if err !=nil{
// 		fmt.Printf("Error: %s \n",err)
// 		return
// 	}
// 	if string(owner_exists) == "false"{
// 		fmt.Printf("Owner does not exist! Owner has to be registered.\n")
// 		return
// 	}

// 	owner_valid, err := ownercontract.EvaluateTransaction("IsOwnerActive", ownername)
// 		if err !=nil{
// 			fmt.Printf("Error: %s \n",err)
// 			return
// 		}
// 	owner_valid_str:=string(owner_valid)
// 	fmt.Printf("OWNER VALID VALUE= %s", owner_valid_str)

// 	if owner_valid_str =="false"{
// 		fmt.Printf("%s",owner_valid_str)
// 		fmt.Printf("Owner %s is not active! Registration of owner is required \n", ownername)
// 		return
// 	}
// 	// GETTING OWNER ID 
// 	ownerIDextract, err:= ownercontract.EvaluateTransaction("ReturnOwnerID", ownername)
// 	if err !=nil{
// 		fmt.Errorf("Error: %s",err)
// 		return
// 	}
// 	ownerID := string(ownerIDextract)
	
	
// 		fmt.Printf("\n--> Creating Asset : %s\n", name)
// 		result, err := contract.SubmitTransaction("CreateAsset", name) // SubmitTransaction returns results of a transaction only after its commited
// 		fmt.Printf("\n--> Submit Transaction Returned : %s , %s\n", string(result), err)
// 		return	
		
//  }


func increaseValue(contract *client.Contract, name string, incVal string) {

	// GETTING OWNER ID 
	// ownerIDextract, err:= ownercontract.EvaluateTransaction("ReturnOwnerID", ownername)
	// if err !=nil{
	// 	fmt.Errorf("Error: %s",err)
	// 	return
	// }
	// ownerID := string(ownerIDextract)

	fmt.Printf("Name : %s , IncreaseValue: %s ", name, incVal)

	evaluatedAsset, err := contract.SubmitTransaction("IncreaseAsset", name, incVal)
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(evaluatedAsset), err)
}

func decreaseValue(contract *client.Contract, name string, decVal string) {


	// GETTING OWNER ID 
	// ownerIDextract, err:= ownercontract.EvaluateTransaction("ReturnOwnerID", ownername)
	// if err !=nil{
	// 	fmt.Errorf("Error: %s",err)
	// 	return
	// }
	// ownerID := string(ownerIDextract)

	fmt.Printf("Name : %s , DecreaseValue: %s ", name, decVal)

	evaluatedAsset, err := contract.SubmitTransaction("DecreaseAsset", name, decVal)
	fmt.Printf("\n------> After SubmitTransaction:%s , %s \n", string(evaluatedAsset), err)
}




func transferAsset(contract *client.Contract,  name string) {

	// // verifying owner existence
	// owner_exists, err := ownercontract.EvaluateTransaction("OwnerExistence", owner_name)
	// if err !=nil{
	// 	fmt.Printf("Error: %s \n",err)
	// 	return
	// }
	// new_owner_exists, err := ownercontract.EvaluateTransaction("OwnerExistence", new_owner_name)
	// if err !=nil{
	// 	fmt.Printf("Error: %s \n",err)
	// 	return
	// }

	// if string(owner_exists) == "false"{
	// 	fmt.Printf("Owner %s does not exist! Owner has to be registered.\n", owner_name )
	// 	return
	// }
	// if string(new_owner_exists) == "false" {
	// 	fmt.Printf("Owner %s does not exist! Owner has to be registered.\n", new_owner_name)
	// 	return
	// }
	
	// if owners exist---->

	// GETTING OWNER IDs : ReturnOwnerID() 

	// newownerIDextract, err:= ownercontract.EvaluateTransaction("ReturnOwnerID", new_owner_name)
	// if err !=nil{
	// 	fmt.Errorf("Error: %s",err)
	// 	return
	// }
	// newownerID := string(newownerIDextract)

	fmt.Printf("Asset name : %s , Transfer asset to: %s ", name, user)

	evaluatedAsset, err := contract.SubmitTransaction("TransferAsset", name)
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

func deleteOwner(contract *client.Contract , name string){


	_, err := contract.SubmitTransaction("DeleteOwner", name) 
	fmt.Printf("\n------> After SubmitTransaction: %s \n",  err)
}


func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return "error in parsing JSON"
	}
	return prettyJSON.String()
}
