#!/bin/bash

#Usage: 5A_AddOrdererDockerNode.sh <physical node no> <orderer_port> 


. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 

  if [[ $# -lt 4 ]] ; then
    fatalln "Usage: 5A_AddOrdererDockerNode.sh <Orderer node Number:1,2 etc.> <orderer_port> <orderer_admin_port> <orderer_operation_port> [--no-reg]"
  fi

  PHYSICAL_SVR_NODE=$1
  DOCKER_ORDERER_PORT=$2
  DOCKER_ORDERER_ADMIN_PORT=$3
  DOCKER_ORDERER_OPS_PORT=$4
  NO_REG=$5

setCANode 1

. ./RegisterEnroll.sh

function createDockerCompose()
{

  #First Register the user
  if [[ "$NO_REG" == "--no-reg" ]]; then
     warnln "Skipping registering user ${ORDERER_USER}."
  else 
     RegisterUser ${ORDERER_CA_HOME} ${ORDERER_USER} ${ORDERER_PASSWORD} orderer
  fi
  
  ORDERER_NODE_FOLDER=organizations/ordererOrganizations/${DOMAIN}
  
  ORDERER_NODE_HOME=${TALLY_HOME}/${ORDERER_NODE_FOLDER}
  
  /bin/rm -rf "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}"
  
  export FABRIC_CA_CLIENT_HOME=${ORDERER_CA_HOME}/client

  infoln "Generating the orderer msp"

  fabric-ca-client enroll -u https://${ORDERER_USER}:${ORDERER_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT} --caname ${ORDERER_CA_NAME} --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=orderer -M "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/msp" --csr.hosts ${ORDERER_HOST}.${DOMAIN} --tls.certfiles "${ORDERER_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll orderer msp."
  
  cp "${ORDERER_NODE_HOME}/msp/config.yaml" "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  fabric-ca-client enroll -u https://${TLS_ADMIN_USER}:${TLS_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT} --caname ${TLS_CA_NAME} -M "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls" --enrollment.profile tls --csr.hosts ${ORDERER_HOST}.${DOMAIN} --tls.certfiles "${TLS_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll orderer tls."

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/tlscacerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/ca.crt"
  verifyResult $? "Unable to copy tls certificate."
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/signcerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/server.crt"
  verifyResult $? "Unable to copy tls certificate."
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/keystore/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/server.key"
  verifyResult $? "Unable to copy tls certificate."

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/msp/tlscacerts"
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/tlscacerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"
  verifyResult $? "Unable to create tlc ca certificate."


  infoln "Deleting remote node folder ..."


  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} /bin/rm -rf ${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}

  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} /bin/mkdir -p ${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers

  scp -C -i ${ORDERER_HOST_KEY} -r ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}:${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}
  verifyResult $? "Unable to copy ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT} to ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}"


  #create orderer config file
	/bin/cp compose/docker/orderer_template.yaml ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml

  ORDERER_HOME_ESCAPED=`echo ${ORDERER_HOME_REMOTE} | sed -e 's/\//\\\\\//g'`
  sed -i "s/\${DOCKER_ORDERER_PORT}/${DOCKER_ORDERER_PORT}/g"             ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_MSPID}/${ORDERER_MSPID}/g"           ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_HOME}/${ORDERER_HOME_ESCAPED}/g"     ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${DOCKER_ORDERER_ADMIN_PORT}/${DOCKER_ORDERER_ADMIN_PORT}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${DOCKER_ORDERER_OPS_PORT}/${DOCKER_ORDERER_OPS_PORT}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_HOST}/${ORDERER_HOST}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${DOMAIN}/${DOMAIN}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_PORT}/${ORDERER_PORT}/g"             ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_ADMIN_PORT}/${ORDERER_ADMIN_PORT}/g"             ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_OPS_PORT}/${ORDERER_OPS_PORT}/g"             ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."


  #transfer the config file to orderer machine
  scp -C -i ${ORDERER_HOST_KEY} ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}:${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml
  verifyResult $? "Unable to copy orderer.yaml to ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}"

  successln "Successfully created orderer docker compose for port ${DOCKER_ORDERER_PORT} on node $1"
}

function DockerContainerId()
{
  echo $(ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} docker ps -f name=${ORDERER_HOST}.${DOMAIN}_${DOCKER_ORDERER_PORT} -q)
}

function startDockerContainer()
{
  infoln "Starting docker container ..."
  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} "docker-compose -f ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/orderer.yaml up -d"
}

function add_orderer_node()
{

  setOrdererNode $PHYSICAL_SVR_NODE

  #update to first peer node
  setPeerNode 1  
  setup_peer_paths

  /bin/rm -rf temp
  mkdir temp
  
  infoln "Fetching config from ${ORDERER_HOST}..."
  peer channel fetch config temp/config_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -c ${CHANNEL_ID}
  verifyResult $? "Fetch config failed"
  
  infoln "Converting config to json ..."
  configtxlator proto_decode --input temp/config_block.pb --type common.Block | jq .data.data[0].payload.data.config > temp/config.json  
  verifyResult $? "Conversion failed"
  
  infoln "Adding new orderer node ${ORDERER_HOST}:${DOCKER_ORDERER_PORT} ..."
  TLS_FILE=${ORDERER_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/server.crt
  echo "{\"client_tls_cert\":\"$(cat $TLS_FILE | base64 -w 0)\",\"host\":\"$ORDERER_HOST.$DOMAIN\",\"port\":$DOCKER_ORDERER_PORT,\"server_tls_cert\":\"$(cat $TLS_FILE | base64 -w 0)\"}" > temp/consenter.json
  jq ".channel_group.groups.Orderer.values.ConsensusType.value.metadata.consenters += [$(cat temp/consenter.json)]" temp/config.json > temp/modified_orderer_config.json
  verifyResult $? "Adding failed"


  infoln "Converting json to protobuf ..."
  configtxlator proto_encode --input temp/config.json --type common.Config --output temp/config.pb
  verifyResult $? "Conversion failed"
  configtxlator proto_encode --input temp/modified_orderer_config.json --type common.Config --output temp/modified_orderer_config.pb
  verifyResult $? "Conversion failed"
  
  infoln "Getting update delta ..."
  configtxlator compute_update --channel_id $CHANNEL_ID --original temp/config.pb --updated temp/modified_orderer_config.pb --output temp/orderer_update.pb  
  verifyResult $? "Updation failed"

  infoln "Creating update json ..."
  configtxlator proto_decode --input temp/orderer_update.pb --type common.ConfigUpdate | jq . > temp/orderer_update.json
  verifyResult $? "Updation failed"

  infoln "Creating payload json ..."
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'${CHANNEL_ID}'", "type":2}},"data":{"config_update":'$(cat temp/orderer_update.json)'}}}' | jq . > temp/orderer_update_in_envelope.json
  
  infoln "Creating payload protobuf ..."
  configtxlator proto_encode --input temp/orderer_update_in_envelope.json --type common.Envelope --output temp/orderer_update_in_envelope.pb
  verifyResult $? "Creation failed"

  infoln "Update config to ${ORDERER_HOST}..."
  
  export CORE_PEER_MSPCONFIGPATH=${ORDERER_HOME}/users/Admin@${DOMAIN}/msp
  export CORE_PEER_LOCALMSPID=${ORDERER_MSPID}
  export CORE_PEER_TLS_ROOTCERT_FILE=${ORDERER_HOME}/orderers/${ORDERER_HOST}_${DOCKER_ORDERER_PORT}/tls/ca.crt

  peer channel update -f temp/orderer_update_in_envelope.pb temp/config_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -c ${CHANNEL_ID} 
  verifyResult $? "Channel update failed"

  /bin/rm -rf temp
  
  setup_peer_paths
  peer channel getinfo -c ${CHANNEL_ID}
  
  successln "New Orderer added to orderer service"
}

function isChannelActive()
{
   
   export OSN_TLS_CA_ROOT_CERT=${TLS_CA_HOME}/ca-cert.pem
   export ADMIN_TLS_SIGN_CERT=${TALLY_HOME}/admin_client/client-tls-cert.pem
   export ADMIN_TLS_PRIVATE_KEY=${TALLY_HOME}/admin_client/client-tls-key.pem
  
   infoln "Checking channel ${CHANNEL_ID} on Orderer $PHYSICAL_SVR_NODE:$DOCKER_ORDERER_PORT ..."
   list_json=$(osnadmin channel list --channelID ${CHANNEL_ID} -o ${ORDERER_HOST}.${DOMAIN}:${DOCKER_ORDERER_ADMIN_PORT} --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY | tail -n +2)
   channel_name=$(echo $list_json | jq -r 'try (.name)')
   channel_status=$(echo $list_json | jq -r 'try (.status)')

   if [[ "$channel_name" == "$CHANNEL_ID" ]] && [[ "$channel_status" == "active" ]]; then
      return 0
   fi

   return 1

}
function  setup_channel()
{
  
  isChannelActive
  if [[ $? -eq 0 ]];then
   warnln "Channel $CHANNEL_ID already active for orderer $PHYSICAL_SVR_NODE:$DOCKER_ORDERER_PORT, skippping ..."
   return
  fi

  OSN_TLS_CA_ROOT_CERT=${TLS_CA_HOME}/ca-cert.pem
  ADMIN_TLS_SIGN_CERT=${TALLY_HOME}/admin_client/client-tls-cert.pem
  ADMIN_TLS_PRIVATE_KEY=${TALLY_HOME}/admin_client/client-tls-key.pem

  /bin/rm -rf temp
  mkdir temp
  
  infoln "Fetching genesis block from ${ORDERER_HOST}..."
  #update to first peer node
  setPeerNode 1  
  setup_peer_paths
  peer channel fetch 0 temp/genesis_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_PORT} --tls --cafile "${ORDERER_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem" -c ${CHANNEL_ID}
  verifyResult $? "Fetch genesis block failed"
  
  
  infoln "Setting up channel $CHANNEL_ID ..."
  osnadmin channel join --channelID ${CHANNEL_ID}  --config-block temp/genesis_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${DOCKER_ORDERER_ADMIN_PORT} --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
  verifyResult $? "Channel setup failed."

  /bin/rm -rf temp

  successln "Orderer $PHYSICAL_SVR_NODE:$DOCKER_ORDERER_PORT : created and joined channel $CHANNEL_ID with genesis block."
  
}

setOrdererDockerNode $PHYSICAL_SVR_NODE $DOCKER_ORDERER_PORT
CONTAINER_ID=$(DockerContainerId)

if [[ "$CONTAINER_ID" != "" ]];then
  fatalln "There is already a container(id=$CONTAINER_ID) setup, please chose different port or use different orderer node or remove this container."
fi

createDockerCompose 

startDockerContainer

CONTAINER_ID=$(DockerContainerId)
if [[ "$CONTAINER_ID" != "" ]];then
  successln "Container(id=$CONTAINER_ID) is setup and started."
else
  fatalln "Failed to setup and start container"
fi


add_orderer_node

setup_channel