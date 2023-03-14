#!/bin/bash

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 
. ./SetPeerNode.sh 

setCANode 1

#Usage: 14A_UpgradeChaincode.sh <ChaincodeName>
function printHelp()
{
  infoln "Usage: 14A_UpgradeChaincode.sh <ChaincodeName> [flags]"
  infoln "Flags:"
  infoln "    -v version     : Version of the chaincode, default: 1.0"
  infoln "    -h             : print this help"
}
if [[ $# -lt 1 ]] ; then
  printHelp  
  exit 1
fi

CC_NAME=$1
shift

#defaults
CC_VERSION="1.0"

while [[ $# -ge 1 ]] ; do
  key="$1"
  case $key in
  -h )
    printHelp 
    exit 0
    ;;
  -v )
    CC_VERSION="$2"
    shift
    ;;
  * )
    echo "Unknown flag: $key"
    printHelp
    exit 1
    ;;
  esac
  shift
done

infoln "Running chaincode upgrade with: "
infoln "Name          = $CC_NAME"
infoln "Version       = $CC_VERSION"


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
    errorln "jq command not found..."
    errorln
    errorln "Follow the instructions in the Fabric docs to install the prereqs"
    errorln "https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html"
    exit 1
  fi
}


# installChaincode PEER ORG
function upgradeChaincode() {

    infoln "Upgrading chaincode ${CC_NAME}..."
    peer chaincode upgrade -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -C ${CHANNEL_ID} -n ${CC_NAME} -v ${CC_VERSION} 
    verifyResult $? "Chaincode Upgradation failed on ${ORDERER_HOST}"
    successln "Chaincode is upgraded on ${ORDERER_HOST}"
}

function upgradeCC()
{
  
   #Use orderer node as 1
  
   setPeerNode 1
   setOrdererNode 1

   setup_peer_paths
   
   #check for prerequisites
   checkPrereqs
   
  
   ## package the chaincode
   upgradeChaincode
   

}

upgradeCC

