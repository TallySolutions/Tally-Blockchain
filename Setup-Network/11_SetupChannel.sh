#!/bin/bash

#Usage: 11_SetupChannel.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 
. ./SetPeerNode.sh 

setCANode 1

function  create_genesis()
{
  
  #create configtx file
  /bin/cp configtx-template.yaml ${TALLY_HOME}/admin_client/configtx.yaml
  
  ORDERER_HOME_ESCAPED=`echo ${ORDERER_HOME} | sed -e 's/\//\\\\\//g'`
  PEER_HOME_ESCAPED=`echo ${PEER_HOME} | sed -e 's/\//\\\\\//g'`
  sed -i "s/\${ORDERER_MSPID}/${ORDERER_MSPID}/g"             ${TALLY_HOME}/admin_client/configtx.yaml 
  verifyResult $? "Update config failed."
  sed -i "s/\${ORDERER_HOME}/${ORDERER_HOME_ESCAPED}/g"       ${TALLY_HOME}/admin_client/configtx.yaml 
  verifyResult $? "Update config failed."
  sed -i "s/\${ORDERER_HOST_PREFIX}/${ORDERER_HOST_PREFIX}/g" ${TALLY_HOME}/admin_client/configtx.yaml 
  verifyResult $? "Update config failed."
  sed -i "s/\${DOMAIN}/${DOMAIN}/g"                           ${TALLY_HOME}/admin_client/configtx.yaml 
  verifyResult $? "Update config failed."
  sed -i "s/\${ORDERER_PORT}/${ORDERER_PORT}/g"               ${TALLY_HOME}/admin_client/configtx.yaml 
  verifyResult $? "Update config failed."
  sed -i "s/\${PEER_MSPID}/${PEER_MSPID}/g"                   ${TALLY_HOME}/admin_client/configtx.yaml 
  verifyResult $? "Update config failed."
  sed -i "s/\${PEER_HOME}/${PEER_HOME_ESCAPED}/g"             ${TALLY_HOME}/admin_client/configtx.yaml 
  verifyResult $? "Update config failed."

  export FABRIC_CFG_PATH=${TALLY_HOME}/admin_client/
  
  mkdir -p ${TALLY_HOME}/admin_client/${CHANNEL_ID}
  
  infoln "Creating Channel genesis block ..."
  
  configtxgen -profile TallyApplicationGenesis -outputBlock ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb -channelID ${CHANNEL_ID}
  verifyResult $? "Genesis block creation failed."

  successln "Genesis block created."

}
function isChannelActive()
{
   if [[ $# -lt 1 ]] ; then
     fatalln "Usage: isChannelActive <Order node Number:1,2 etc.>"
   fi
   setOrdererNode $1 
   
   export OSN_TLS_CA_ROOT_CERT=${TLS_CA_HOME}/ca-cert.pem
   export ADMIN_TLS_SIGN_CERT=${TALLY_HOME}/admin_client/client-tls-cert.pem
   export ADMIN_TLS_PRIVATE_KEY=${TALLY_HOME}/admin_client/client-tls-key.pem
  
   infoln "Checking channel ${CHANNEL_ID} on Orderer $1 ..."
   list_json=$(osnadmin channel list --channelID ${CHANNEL_ID} -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_ADMIN_PORT} --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY | tail -n +2)
   channel_name=$(echo $list_json | jq -r 'try (.name)')
   channel_status=$(echo $list_json | jq -r 'try (.status)')

   if [[ "$channel_name" == "$CHANNEL_ID" ]] && [[ "$channel_status" == "active" ]]; then
      return 0
   fi

   return 1

}
function  setup_channel()
{
  
  if [[ $# -lt 1 ]] ; then
    fatalln "Usage: setup_channel <Order node Number:1,2 etc.>"
  fi

  setOrdererNode $1 

  isChannelActive $1
  if [[ $? -eq 0 ]];then
   warnln "Channel $CHANNEL_ID already active for orderer $1, skippping ..."
   return
  fi

  OSN_TLS_CA_ROOT_CERT=${TLS_CA_HOME}/ca-cert.pem
  ADMIN_TLS_SIGN_CERT=${TALLY_HOME}/admin_client/client-tls-cert.pem
  ADMIN_TLS_PRIVATE_KEY=${TALLY_HOME}/admin_client/client-tls-key.pem
  
  infoln "Setting up channel $CHANNEL_ID ..."
  osnadmin channel join --channelID ${CHANNEL_ID}  --config-block ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_ADMIN_PORT} --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
  verifyResult $? "Channel setup failed."

  successln "Orderer $1 : created and joined channel $CHANNEL_ID with genesis block."
  
}

function setup_peer_paths()
{
  export FABRIC_CFG_PATH=${PEER_HOME}/peers/${PEER_HOST}
  export CORE_PEER_MSPCONFIGPATH=${PEER_HOME}/users/Admin@${DOMAIN}/msp
  export CORE_PEER_ADDRESS=${PEER_HOST}.${DOMAIN}:${PEER_PORT}
  export CORE_PEER_LOCALMSPID=${PEER_MSPID}
  export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_HOME}/peers/${PEER_HOST}/tls/ca.crt
}

function join_peer()
{
    if [[ $# -lt 1 ]] ; then
      fatalln "Usage: join_peer <Peer node Number:1,2 etc.>"
    fi
  
  setPeerNode $1
  setup_peer_paths
  
  infoln "Joining Channel ${CHANNEL_ID} for peer $1..."
  
  peer channel getinfo --channelID ${CHANNEL_ID}
  if [[ $? -eq 0 ]];then
   warnln "Channel $CHANNEL_ID already exists for peer $1, skippping ..."
   return
  fi

  peer channel join -b ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb
  verifyResult $? "Channel join failed."

  successln "Peer $1 joined channel $CHANNEL_ID"
  

}



create_genesis

setup_channel 1
setup_channel 2

join_peer 1
join_peer 2
join_peer 3

