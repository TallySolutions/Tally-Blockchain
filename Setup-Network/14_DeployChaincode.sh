#!/bin/bash

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 
. ./SetPeerNode.sh 

setCANode 1

#Usage: 14_DeployChaincode.sh <ChaincodeName> <ChainCodePath>
function printHelp()
{
  infoln "Usage: 11_DeployChaincode.sh <ChaincodeName> <ChainCodePath> [flags]"
  infoln "Flags:"
  infoln "    -v version     : Version of the chaincode, default: 1.0"
  infoln "    -s int         : The sequence number of the chaincode definition for the channel, default: 1"
  infoln "    -f function    : Init function to invoked after deploying chaincode, default: NA"
  infoln "    -d int         : delay in seconds, before retry, default: 3"
  infoln "    -r int         : No of retries, default: 5"
  infoln "    -h             : print this help"
}
if [[ $# -lt 2 ]] ; then
  printHelp  
  exit 1
fi

CC_NAME=$1
shift
CC_SRC_PATH=$1
shift

#defailts
CC_VERSION="1.0"
CC_SEQUENCE=1
CC_INIT_FCN="NA"
DELAY=3
MAX_RETRY=5

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
  -s )
    CC_SEQUENCE="$2"
    shift
    ;;
  -r )
    MAX_RETRY="$2"
    shift
    ;;
  -d )
    DELAY="$2"
    shift
    ;;
  -f )
    CC_INIT_FCN="$2"
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

infoln "Running chaincode deploy with: "
infoln "Name          = $CC_NAME"
infoln "Path          = $CC_SRC_PATH"
infoln "Version       = $CC_VERSION"
infoln "Sequence      = $CC_SEQUENCE"
infoln "Init Function = $CC_INIT_FCN"
infoln "Delay         = $DELAY"
infoln "Max Retry     = $MAX_RETRY"


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
  successln "Chaincode is packaged"
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
  successln "Chaincode is installed on ${PEER_HOST}"
}

# queryInstalled PEER ORG
function queryInstalled() {
  set -x
  peer lifecycle chaincode queryinstalled --output json | jq -r 'try (.installed_chaincodes[].package_id)' | grep ^${PACKAGE_ID}$ >&${CC_PKG_PATH}/log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat ${CC_PKG_PATH}/log.txt
  verifyResult $res "Query installed on ${PEER_HOST} has failed"
  successln "Query installed successful on ${PEER_HOST} on channel"
}

# approveForMyOrg VERSION PEER ORG
function approveForTally() {
  set -x
  peer lifecycle chaincode approveformyorg -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" --channelID ${CHANNEL_ID} --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&${CC_PKG_PATH}/log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat ${CC_PKG_PATH}/log.txt
  verifyResult $res "Chaincode definition approved on ${PEER_HOST} on channel '$CHANNEL_ID' failed"
  successln "Chaincode definition approved on ${PEER_HOST} on channel '$CHANNEL_ID'"
}

function checkCommitReadiness()
{
  infoln "Checking the commit readiness of the chaincode definition on $PEER_HOST on channel '$CHANNEL_ID'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    infoln "Attempting to check the commit readiness of the chaincode definition on $PEER_HOST, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_ID --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --output json >&${CC_PKG_PATH}/log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=0
    for var in "$@"; do
      grep "$var" ${CC_PKG_PATH}/log.txt &>/dev/null || let rc=1
    done
    COUNTER=$(expr $COUNTER + 1)
  done
  cat ${CC_PKG_PATH}/log.txt
  if test $rc -eq 0; then
    successln "Checking the commit readiness of the chaincode definition successful on $PEER_HOST on channel '$CHANNEL_ID'"
  else
    fatalln "After $MAX_RETRY attempts, Check commit readiness result on $PEER_HOST is INVALID!"
  fi
}

function commitChaincodeDefinition()
{
  # while 'peer chaincode' command can get the orderer endpoint from the
  # peer (if join was successful), let's supply it directly as we know
  # it using the "-o" option
  set -x
  peer lifecycle chaincode commit -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" --channelID $CHANNEL_ID --name ${CC_NAME} --peerAddresses ${PEER_HOST}.${DOMAIN}:${PEER_PORT} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --tlsRootCertFiles ${PEER_HOME}/peers/${PEER_HOST}/tls/ca.crt >&${CC_PKG_PATH}/log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat ${CC_PKG_PATH}/log.txt
  verifyResult $res "Chaincode definition commit failed on $PEER_HOST on channel '$CHANNEL_ID' failed"
  successln "Chaincode definition committed on channel '$CHANNEL_ID'"

}

function queryCommitted()
{
  EXPECTED_RESULT="Version: ${CC_VERSION}, Sequence: ${CC_SEQUENCE}, Endorsement Plugin: escc, Validation Plugin: vscc"
  infoln "Querying chaincode definition on peer0.org${ORG} on channel '$CHANNEL_ID'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    echo "Attempting to Query committed status on ${PEER_HOST}, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_ID --name ${CC_NAME} >&${CC_PKG_PATH}/log.txt
    res=$?
    { set +x; } 2>/dev/null
    test $res -eq 0 && VALUE=$(cat $CC_PKG_PATH/log.txt | grep -o '^Version: '$CC_VERSION', Sequence: [0-9]*, Endorsement Plugin: escc, Validation Plugin: vscc')
    test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
    COUNTER=$(expr $COUNTER + 1)
  done
  cat $CC_PKG_PATH/log.txt
  if test $rc -eq 0; then
    successln "Query chaincode definition successful on ${PEER_HOST} on channel '$CHANNEL_NAME'"
  else
    fatalln "After $MAX_RETRY attempts, Query chaincode definition result on ${PEER_HOST} is INVALID!"
  fi
}

function chaincodeInvokeInit() {

  # while 'peer chaincode' command can get the orderer endpoint from the
  # peer (if join was successful), let's supply it directly as we know
  # it using the "-o" option
  set -x
  fcn_call='{"function":"'${CC_INIT_FCN}'","Args":[]}'
  infoln "invoke fcn call:${fcn_call}"
  peer chaincode invoke -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -C $CHANNEL_ID -n ${CC_NAME} --peerAddresses ${PEER_HOST}.${DOMAIN}:${PEER_PORT}  --isInit -c ${fcn_call} >&$CC_PKG_PATH/log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat $CC_PKG_PATH/log.txt
  verifyResult $res "Invoke execution on $PEER_HOST failed "
  successln "Invoke transaction successful on $PEER_HOST on channel '$CHANNEL_NAME'"
}

function deployCC()
{
  
   #Use peer host 1
  
   setPeerNode 1

   CC_PKG_PATH=${TALLY_HOME}/admin_client/chaincode/${CC_NAME}
   
   setup_peer_paths
   
   #check for prerequisites
   checkPrereqs
   
   
   ## package the chaincode
   infoln "Packaging chaincode on ${PEER_HOST}"
   packageChaincode
   
   ## Install chaincode on peer0.org1 and peer0.org2
   infoln "Installing chaincode on ${PEER_HOST}"
   installChaincode 

   ## query whether the chaincode is installed
   infoln "Querying chaincode on ${PEER_HOST}"
   queryInstalled

   ## approve the definition 
   approveForTally

   ## check whether the chaincode definition is ready to be committed
   ## expect Tally to have approved 
   checkCommitReadiness "\"Tally\": true" 

   ## now that we know for sure both orgs have approved, commit the definition
   commitChaincodeDefinition

   ## query on both orgs to see that the definition committed successfully
   queryCommitted 

   ## Invoke the chaincode - this does require that the chaincode have the 'initLedger'
   ## method defined
   if [ "$CC_INIT_FCN" = "NA" ]; then
     infoln "Chaincode initialization is not required"
   else
     infoln "Invoking Chaincode ..."
     chaincodeInvokeInit
   fi

   successln "Chaincode deployed."

}

deployCC 