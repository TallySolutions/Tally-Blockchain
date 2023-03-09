#!/bin/bash

#Usage: 9_SetupChannel.sh 

function fatalln() {
  echo "ERROR: $1"
  exit 1
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}

function  create_genesis()
{
  . ./SetGlobalVariables.sh 1
  
  #create configtx file
  /bin/cp configtx-template.yaml ${TALLY_HOME}/admin_client/configtx.yaml
  
  ORDERER_HOME_ESCAPED=`echo ${ORDERER_HOME} | sed -e 's/\//\\\\\//g'`
  PEER_HOME_ESCAPED=`echo ${PEER_HOME} | sed -e 's/\//\\\\\//g'`
  sed -e "s/\${ORDERER_MSPID}/${ORDERER_MSPID}/g"             ${TALLY_HOME}/admin_client/configtx.yaml   > ${TALLY_HOME}/admin_client/configtx.yaml.1 
  sed -e "s/\${ORDERER_HOME}/${ORDERER_HOME_ESCAPED}/g"       ${TALLY_HOME}/admin_client/configtx.yaml.1 > ${TALLY_HOME}/admin_client/configtx.yaml.2 
  sed -e "s/\${ORDERER_HOST_PREFIX}/${ORDERER_HOST_PREFIX}/g" ${TALLY_HOME}/admin_client/configtx.yaml.2 > ${TALLY_HOME}/admin_client/configtx.yaml.3 
  sed -e "s/\${DOMAIN}/${DOMAIN}/g"                           ${TALLY_HOME}/admin_client/configtx.yaml.3 > ${TALLY_HOME}/admin_client/configtx.yaml.4 
  sed -e "s/\${ORDERER_PORT}/${ORDERER_PORT}/g"               ${TALLY_HOME}/admin_client/configtx.yaml.4 > ${TALLY_HOME}/admin_client/configtx.yaml.5 
  sed -e "s/\${PEER_MSPID}/${PEER_MSPID}/g"                   ${TALLY_HOME}/admin_client/configtx.yaml.5 > ${TALLY_HOME}/admin_client/configtx.yaml.6 
  sed -e "s/\${PEER_HOME}/${PEER_HOME_ESCAPED}/g"             ${TALLY_HOME}/admin_client/configtx.yaml.6 > ${TALLY_HOME}/admin_client/configtx.yaml.7 

  /bin/rm ${TALLY_HOME}/admin_client/configtx.yaml ${TALLY_HOME}/admin_client/configtx.yaml.1 ${TALLY_HOME}/admin_client/configtx.yaml.2 ${TALLY_HOME}/admin_client/configtx.yaml.3 ${TALLY_HOME}/admin_client/configtx.yaml.4 ${TALLY_HOME}/admin_client/configtx.yaml.5 ${TALLY_HOME}/admin_client/configtx.yaml.6
  /bin/mv ${TALLY_HOME}/admin_client/configtx.yaml.7  ${TALLY_HOME}/admin_client/configtx.yaml

  export FABRIC_CFG_PATH=${TALLY_HOME}/admin_client/
  
  mkdir -p ${TALLY_HOME}/admin_client/${CHANNEL_ID}
  
  echo "Creating Channel genesis block ..."
 
  
  configtxgen -profile TallyApplicationGenesis -outputBlock ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb -channelID ${CHANNEL_ID}

  configtxlator proto_decode --input ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb --type common.Block | jq .data.data[0].payload.data.config > ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.json

  echo "Genesis block created ====>"

  cat   ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.json

}

function  setup_channel()
{
  
  if [[ $# -lt 1 ]] ; then
    echo "Usage: setup_channel <Order node Number:1,2 etc.>"
    exit 1
  fi
  
  . ./SetGlobalVariables.sh $1

  OSN_TLS_CA_ROOT_CERT=${TLS_CA_HOME}/ca-cert.pem
  ADMIN_TLS_SIGN_CERT=${TALLY_HOME}/admin_client/client-tls-cert.pem
  ADMIN_TLS_PRIVATE_KEY=${TALLY_HOME}/admin_client/client-tls-key.pem
  
  set -x

  osnadmin channel join --channelID ${CHANNEL_ID}  --config-block ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_ADMIN_PORT} --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
  
  { set +x; } 2> /dev/null
  
  
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
      echo "Usage: join_peer <Peer node Number:1,2 etc.>"
      exit 1
    fi
  
  . ./SetGlobalVariables.sh $1

  setup_peer_paths
  
  echo "Joining Channel ${CHANNEL_ID} ..."
  
  set -x
  peer channel join -b ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb
  { set +x; } 2> /dev/null


}

function set_anchor_peer()
{
    if [[ $# -lt 1 ]] ; then
      echo "Usage: set_anchor_peer <Peer node Number:1,2 etc.>"
      exit 1
    fi
  
  . ./SetGlobalVariables.sh $1

  setup_peer_paths

  /bin/rm -rf temp
  mkdir temp
  
  echo "Fetching config from ${ORDERER_HOST}..."
  peer channel fetch config temp/config_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -c ${CHANNEL_ID}
  verifyResult $? "Fetch config failed"
  
  echo "Converting config to json ..."
  configtxlator proto_decode --input temp/config_block.pb --type common.Block | jq .data.data[0].payload.data.config > temp/config.json  
  verifyResult $? "Conversion failed"
  
  echo "Adding anchor peer ${PEER_HOST} ..."
  jq '.channel_group.groups.Application.groups.'${PEER_MSPID}'.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "'${PEER_HOST}'.'${DOMAIN}'","port": '${PEER_PORT}'}]},"version": "0"}}' temp/config.json > temp/modified_anchor_config.json
  verifyResult $? "Adding failed"
  
  cat temp/modified_anchor_config.json

  echo "Converting json to protobuf ..."
  configtxlator proto_encode --input temp/config.json --type common.Config --output temp/config.pb
  verifyResult $? "Conversion failed"
  configtxlator proto_encode --input temp/modified_anchor_config.json --type common.Config --output temp/modified_anchor_config.pb
  verifyResult $? "Conversion failed"
  
  echo "Getting update delta ..."
  configtxlator compute_update --channel_id $CHANNEL_ID --original temp/config.pb --updated temp/modified_anchor_config.pb --output temp/anchor_update.pb  
  verifyResult $? "Updation failed"

  echo "Creating update json ..."
  configtxlator proto_decode --input temp/anchor_update.pb --type common.ConfigUpdate | jq . > temp/anchor_update.json
  verifyResult $? "Updation failed"

  echo "Creating payload json ..."
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'${CHANNEL_ID}'", "type":2}},"data":{"config_update":'$(cat temp/anchor_update.json)'}}}' | jq . > temp/anchor_update_in_envelope.json
  cat temp/anchor_update_in_envelope.json

  echo "Creating payload protobuf ..."
  configtxlator proto_encode --input temp/anchor_update_in_envelope.json --type common.Envelope --output temp/anchor_update_in_envelope.pb
  verifyResult $? "Creation failed"

  echo "Update config to ${ORDERER_HOST}..."
  peer channel update -f temp/anchor_update_in_envelope.pb temp/config_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -c ${CHANNEL_ID} 
  verifyResult $? "Channel update failed"

  /bin/rm -rf temp
  
  peer channel getinfo -c ${CHANNEL_ID}
  
}

create_genesis

setup_channel 1
setup_channel 2

join_peer 1
join_peer 2
join_peer 3

echo "Waiting 15 seconds ..."
sleep 15

set_anchor_peer 1
set_anchor_peer 2
