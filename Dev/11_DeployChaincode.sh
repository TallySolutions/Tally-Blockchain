#!/bin/bash

#Usage: 11_DeployChaincode.sh <ChaincodeName> <ChainCodePath>

if [[ $# -lt 2 ]] ; then
  echo "Usage: 11_DeployChaincode.sh <ChaincodeName> <ChainCodePath>"
  exit 1
fi

function setup_peer_paths()
{
  export FABRIC_CFG_PATH=${PEER_HOME}/peers/${PEER_HOST}
  export CORE_PEER_MSPCONFIGPATH=${PEER_HOME}/users/Admin@${DOMAIN}/msp
  export CORE_PEER_ADDRESS=${PEER_HOST}.${DOMAIN}:${PEER_PORT}
  export CORE_PEER_LOCALMSPID=${PEER_MSPID}
  export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_HOME}/peers/${PEER_HOST}/tls/ca.crt
}

function checkPrereqs() {
  jq --version > /dev/null 2>&1

  if [[ $? -ne 0 ]]; then
    echo "jq command not found..."
    echo
    echo "Follow the instructions in the Fabric docs to install the prereqs"
    echo "https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html"
    exit 1
  fi
}

function fatalln() {
  echo "ERROR: $1"
  exit 1
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}

packageChaincode() {
  set -x
  /bin/rm -rf ${CC_PKG_PATH}
  /bin/mkdir -p ${CC_PKG_PATH}
  
  peer lifecycle chaincode package ${CC_PKG_PATH}/${CC_NAME}.tar.gz --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} --label ${CC_NAME}_${CC_VERSION} >&${CC_PKG_PATH}/log.txt
  res=$?
  PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_PKG_PATH}/${CC_NAME}.tar.gz)
  { set +x; } 2>/dev/null
  cat ${CC_PKG_PATH}/log.txt
  verifyResult $res "Chaincode packaging has failed"
  echo "Chaincode is packaged"
}

# installChaincode PEER ORG
function installChaincode() {

  set -x
  peer lifecycle chaincode queryinstalled --output json | jq -r 'try (.installed_chaincodes[].package_id)' | grep ^${PACKAGE_ID}$ >&${CC_PKG_PATH}/log.txt
  if test $? -ne 0; then
    peer lifecycle chaincode install ${CC_PKG_PATH}/${CC_NAME}.tar.gz >&${CC_PKG_PATH}/log.txt
    res=$?
  fi
  { set +x; } 2>/dev/null
  cat ${CC_PKG_PATH}/log.txt
  verifyResult $res "Chaincode installation on ${PEER_HOST} has failed"
  echo "Chaincode is installed on ${PEER_HOST}"
}

# queryInstalled PEER ORG
function queryInstalled() {
  set -x
  peer lifecycle chaincode queryinstalled --output json | jq -r 'try (.installed_chaincodes[].package_id)' | grep ^${PACKAGE_ID}$ >&${CC_PKG_PATH}/log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat ${CC_PKG_PATH}/log.txt
  verifyResult $res "Query installed on ${PEER_HOST} has failed"
  echo "Query installed successful on ${PEER_HOST} on channel"
}

# approveForMyOrg VERSION PEER ORG
function approveForTally() {
  set -x
  peer lifecycle chaincode approveformyorg -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" --channelID ${CHANNEL_ID} --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&${CC_PKG_PATH}/log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat ${CC_PKG_PATH}/log.txt
  verifyResult $res "Chaincode definition approved on ${PEER_HOST} on channel '$CHANNEL_ID' failed"
  echo "Chaincode definition approved on ${PEER_HOST} on channel '$CHANNEL_ID'"
}

function deployCC()
{
   if [[ $# -lt 3 ]] ; then
     echo "Usage: 11_DeployChaincode.sh <target peer node no: 1,2 etc.> <ChaincodeName> <ChainCodePath>"
     exit 1
   fi
   
   . ./SetGlobalVariables.sh $1
   
   CC_NAME=$2
   CC_SRC_PATH=$3
   CC_VERSION="1.0"
   CC_PKG_PATH=${TALLY_HOME}/admin_client/chaincode/${CC_NAME}
   
   setup_peer_paths
   
   #check for prerequisites
   checkPrereqs
   
   
   ## package the chaincode
   echo "Packaging chaincode on ${PEER_HOST}"
   packageChaincode
   
   ## Install chaincode on peer0.org1 and peer0.org2
   echo "Installing chaincode on ${PEER_HOST}"
   installChaincode 

   ## query whether the chaincode is installed
   echo "Querying chaincode on ${PEER_HOST}"
   queryInstalled

   ## approve the definition 
   approveForTally

## check whether the chaincode definition is ready to be committed
## expect org1 to have approved and org2 not to
#checkCommitReadiness 1 "\"Org1MSP\": true" "\"Org2MSP\": false"
#checkCommitReadiness 2 "\"Org1MSP\": true" "\"Org2MSP\": false"

## now that we know for sure both orgs have approved, commit the definition
#commitChaincodeDefinition 1 2

## query on both orgs to see that the definition committed successfully
#queryCommitted 1

## Invoke the chaincode - this does require that the chaincode have the 'initLedger'
## method defined
if [ "$CC_INIT_FCN" = "NA" ]; then
  infoln "Chaincode initialization is not required"
else
  echo "Invoking Chaincode ..."
  #chaincodeInvokeInit 1 2
fi

}

deployCC 1 $1 $2