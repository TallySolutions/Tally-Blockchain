# Integer Key Application

A blockchain-based application that implements a mechanism to mainly increase, decrease and transfer variables (i.e. assets- in blockchain terminology). 
The assets in this application have Name, Value, OwnerID, TransferStatus, RequestingUser as their properties.

This application allows users to implement functionalities such as creating an asset, increasing as well as decreasing the value of the asset- based on certain conditions. 
Along with this, users are allowed to transfer ownership of the asset to another user. The users are allowed access to these functionalities based on certain attributes (Attribute Based Access Control aka ABAC).

All the calls to the aforementioned functionalities are implemented as functions in the integerkey chaincode as functions. 
Each call to these functions count as blockchain transactions- the changes made by these transactions are reflected on the ledger.

application-with-owner-cli contains code that performs the calls to the chaincode functions. It works through command line instructions.

Along with the mentioned functionalities, there are some other functionalities such as, reading asset, getting all assets, requesting of transfer, approval of transfer. 

Before running the application, the network has to be set up- instructions for this and the codes involved are in the Tally-Blockchain/SetupNetwork directory.

## SetupNetwork

In order to run the application, one must go to the command line and navigate to the SetupNetwork 

The first step is to Start the Orderer servers. 
An orderer receives a transaction after it has been sent to the network and has been peer-validated. 
It then produces blocks of transactions in a specified order and sends them to the peers who are committing them to the ledger for final approval(commiting peers).

Run the following command to Start the orderer servers.

`./6A_StartOrdererServers.sh`

After this, the peer servers have to be started. In order to do that, run the following command.

`./9A_StartPeerServers.sh`

Once the Orderer and Peer servers are up and running, the integerkeychaincode has to be deployed.

The following code deploys the integerkeychaincode on the existing blockchain channel.

`./14_DeployChaincode.sh <channel_name> <chaincode_name> <path_to_chaincode> -s <sequence_number>`

For example,

`./14_DeployChaincode.sh tally integerkey ../integerkey/chaincodes/integerkey/ -s 40`

