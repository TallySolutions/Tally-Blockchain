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

### Setting up servers and deploying the chaincode

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

Syntax: `./14_DeployChaincode.sh <channel_name> <chaincode_name> <path_to_chaincode> -s <sequence_number>`

For example,

`./14_DeployChaincode.sh tally integerkey ../integerkey/chaincodes/integerkey/ -s 40`

### Creation of Users

This section is also executed under the Setup-Network directory. The users created through the following steps can be set as owners to the assets.

This mechanism of having to register users through the fabric ca client allows us to perform ABAC- Attribute Based Access Control. This is a mechanism that enables users to access certain methods of the chaincode according to the attributes that are part of them. The attributes can include information on the role (approver or creator for transactions), etc. The role can be checked and the attribute based conditions can be mentioned in the chaincode using the following method.

```ctx.GetClientIdentity().AssertAttributeValue("creator", "true")```

This checks whether the user has the attribute "creator" set as "true". The variable this is assigned to can be used in the conditional statements to check whether that user has any authority there.

Here is the command to register a user.

Syntax: ```./17_CreateTallyClientUser.sh <userName> <approverAttributeValue> <creatorAttributeValue>```

```./17_CreateTallyClientUser.sh userA false true```

This registers a user named "userA" and sets the approver attribute to false and the creator value to true.

After registration, a password is generated which should be used to enroll the user.

Syntax: ```./18_EnrollTallyUser.sh <userName> <password generated after registering>```

```./18_EnrollTallyUser.sh userA HXXxzrPFsWFk```

After registration and enrolling the users, the application can now be run as any of the existing users.


## Running the application-with-owner-cli

After the network has been set up and the chaincode has been deployed to the channel, the application can be used to call chaincode functions i.e. evaluate or submit transactions to the ledger- depending upon the function called.

Below are the functionalities that the application provides as well as the syntax for usage.

### Create New Asset

This functionality is used for creating an asset. The user who wants to create the asset is specified as an environment variable along with the go command for asset creation.

The following command creates an asset on the ledger, along with the owner. 

Syntax: ```userid=<userName> go run . <peerNodeName> new <assetName>```

```userid=user4 go run . tbchlfdevpeer01 new asset1```

In the above command, "user4" is the owner, tbchlfdevpeer01 is the name of the peer node, "new" is the parameter which is used to call relevant function in the program, "asset1" is the name of the asset to be created. It calls the createAsset() function which in turn, submits a CreateAsset transaction from the chaincode.

### Read an Asset

This functionality is to read an asset and view its properties.

The below command lets us read asset1.

Syntax: ``` userid=<userName> go run . <peerNodeName> read <assetName> ```

``` userid=user4 go run . tbchlfdevpeer01 read asset1 ```

In the above command, tbchlfdevpeer01 is the name of the peer node, "read" is the parameter which is used to call the relevant function in the program, "asset1" is the name of the asset to be read. It calls the readAsset() function which in turn, evaluates a ReadAsset transaction from the chaincode.

### Increment the Value of an asset

This functionality allows the user to increment the value of an asset by a certain amount. The value is not allowed to go above 20.

The following command is used to increment value of asset1.

Syntax: ```userid=<userName> go run . <peerNodeName> inc <assetName> <incrementValue>```

```userid=user4 go run . tbchlfdevpeer01 inc asset1 5```

This command increments the value of asset1 by 5.

### Decrement the Value of an asset

This functionality allows the user to decrement the value of an asset by a certain amount. The value is not allowed to go below 0.

The following command is used to decrement value of asset1.

Syntax: ```userid=<userName> go run . <peerNodeName> dec <assetName> <decrementValue>```

```userid=user4 go run . tbchlfdevpeer01 dec asset1 3```

This command increments the value of asset1 by 3.

### Requesting Transfer of Asset

If there is another user that wants to transfer an asset owned by another user, the user can request to transfer said asset from the owner user. It is upto the owner user to approve transfer of the requested asset.

The following command allows a user named "user5" to request for transfer of asset1.

Syntax: ```userid=<requestingUserName> go run . <peerNodeName> request_transfer <assetName>```

```userid=user5 go run . tbchlfdevpeer01 request_transfer asset1```

The above command submits a transaction to change the transfer status of asset1 to "requested" (handled in integerkey chaincode). It also sets the name of the requesting user as a property of the asset.

### Approving Transfer of Asset

If the asset is requested, the current owner can choose to approve the transfer of said asset to the requesting user

The following command allows user4 to approve transfer of asset1.

Syntax: ```userid=<approvingUserName> go run . <peerNodeName> approve_transfer <assetName>```

```userid=user4 go run . tbchlfdevpeer01 approve_transfer asset1```

The above command submits a transaction to change the transfer status of asset1 to "approved" (handled in integerkey chaincode).

### Performing transfer of asset

Once transfer of the asset has been approved, it can now be transferred. 

The following command lets user5 transfer asset1 to itself.

Syntax: ```userid=<requestingUserName> go run . <peerNodeName> perform_transfer <assetName>```

```userid=user5 go run . tbchlfdevpeer01 perform_transfer asset1```


### Deletion of an asset

This functionality allows the user to delete an asset. 

The following command alows user6 to delete asset1.

Syntax: ```userid=<userName> go run . <peerNodeName> del <assetName>```

```userid=user5 go run . tbchlfdevpeer01 del asset1```

### Viewing list of all assets

This functionality lets us view all the assets that exist on the ledger.

Syntax: ```userid=<userName> go run . <peerNodeName> list```

```userid=user4 go run . tbchlfdevpeer01 list```

The above command lets user4 view all the assets on the ledger.
