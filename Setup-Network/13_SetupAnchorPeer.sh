#!/bin/bash

#Usage: 11_SetupChannel.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 
. ./SetPeerNode.sh 

setCANode 1


function createAnchorPeerList()
{
  if [[ $# -lt 1 ]] ; then
      fatalln "Usage: createAnchorPeerList <Peer node Number:1,2 etc.> ..."
  fi
  PEER_STRING=""
  for node in $*
  do
    setPeerNode $node
    PEER_STRING=${PEER_STRING}{"host": "'${PEER_HOST}'.'${DOMAIN}'","port": '${PEER_PORT}'},
  done
  echo ${PEER_STRING::-1}
}

function set_anchor_peer()
{
    if [[ $# -lt 1 ]] ; then
      fatalln "Usage: set_anchor_peer <Peer node Number:1,2 etc.> ..."
    fi
  
  #use orderer node 1 to fetch config
  setOrdererNode 1

  PEER_STRING=$(createAnchorPeerList $*)

  #update to first anchor peer node
  setPeerNode $1  
  setup_peer_paths

  /bin/rm -rf temp
  mkdir temp
  
  infoln "Fetching config from ${ORDERER_HOST}..."
  peer channel fetch config temp/config_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -c ${CHANNEL_ID}
  verifyResult $? "Fetch config failed"
  
  infoln "Converting config to json ..."
  configtxlator proto_decode --input temp/config_block.pb --type common.Block | jq .data.data[0].payload.data.config > temp/config.json  
  verifyResult $? "Conversion failed"
  
  infoln "Adding anchor peer ${PEER_HOST} ..."
  jq '.channel_group.groups.Application.groups.'${PEER_MSPID}'.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": ['${PEER_STRING}']},"version": "0"}}' temp/config.json > temp/modified_anchor_config.json
  verifyResult $? "Adding failed"


  infoln "Converting json to protobuf ..."
  configtxlator proto_encode --input temp/config.json --type common.Config --output temp/config.pb
  verifyResult $? "Conversion failed"
  configtxlator proto_encode --input temp/modified_anchor_config.json --type common.Config --output temp/modified_anchor_config.pb
  verifyResult $? "Conversion failed"
  
  infoln "Getting update delta ..."
  configtxlator compute_update --channel_id $CHANNEL_ID --original temp/config.pb --updated temp/modified_anchor_config.pb --output temp/anchor_update.pb  
  verifyResult $? "Updation failed"

  infoln "Creating update json ..."
  configtxlator proto_decode --input temp/anchor_update.pb --type common.ConfigUpdate | jq . > temp/anchor_update.json
  verifyResult $? "Updation failed"

  infoln "Creating payload json ..."
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'${CHANNEL_ID}'", "type":2}},"data":{"config_update":'$(cat temp/anchor_update.json)'}}}' | jq . > temp/anchor_update_in_envelope.json
  cat temp/anchor_update_in_envelope.json

  infoln "Creating payload protobuf ..."
  configtxlator proto_encode --input temp/anchor_update_in_envelope.json --type common.Envelope --output temp/anchor_update_in_envelope.pb
  verifyResult $? "Creation failed"

  infoln "Update config to ${ORDERER_HOST}..."
  peer channel update -f temp/anchor_update_in_envelope.pb temp/config_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -c ${CHANNEL_ID} 
  verifyResult $? "Channel update failed"

  /bin/rm -rf temp
  
  peer channel getinfo -c ${CHANNEL_ID}
  
  successln "Peer(S) $* set as anchor peer"
}

set_anchor_peer 1 2
