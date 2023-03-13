#!/bin/bash

#Usage: 8_CreatePeerNodeMSP.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetPeerNode.sh 

setCANode 1

. ./RegisterEnroll.sh

function createPeerNode()
{

  if [[ $# -lt 1 ]] ; then
    fataln "Usage: create <Peer no Number:1,2 etc.>"
  fi

  setPeerNode $1
  
  #First register user
  
  RegisterUser ${TALLY_CA_HOME} ${PEER_USER} ${PEER_PASSWORD} peer
  
  
  /bin/rm -rf "${PEER_NODE_HOME}/peers/${PEER_HOST}"
  
  export FABRIC_CA_CLIENT_HOME=${TALLY_CA_HOME}/client

  infoln "Generating the peer msp"
  fabric-ca-client enroll -u https://${PEER_USER}:${PEER_PASSWORD}@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT} --caname ${TALLY_CA_NAME} --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=client -M "${PEER_NODE_HOME}/peers/${PEER_HOST}/msp" --csr.hosts ${PEER_HOST}.${DOMAIN} --tls.certfiles "${TALLY_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll peer msp."

  cp "${PEER_NODE_HOME}/msp/config.yaml" "${PEER_NODE_HOME}/peers/${PEER_HOST}/msp/config.yaml"
  verifyResult $? "Unable to create peer config."

  echo "Generating the peer-tls certificates"
  fabric-ca-client enroll -u https://${TLS_ADMIN_USER}:${TLS_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT} --caname ${TLS_CA_NAME} -M "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls" --enrollment.profile tls --csr.hosts ${PEER_HOST}.${DOMAIN} --tls.certfiles "${TLS_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll peer tls."

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/tlscacerts/"* "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/ca.crt"
  verifyResult $? "Unable to copy ca.crt."
  cp "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/signcerts/"* "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/server.crt"
  verifyResult $? "Unable to copy server.crt."
  cp "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/keystore/"* "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/server.key"
  verifyResult $? "Unable to copy server.key."

  # Copy peer org's CA cert to peer's /msp/tlscacerts directory (for use in the peer MSP definition)
  mkdir -p "${PEER_NODE_HOME}/peers/${PEER_HOST}/msp/tlscacerts"
  cp "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/tlscacerts/"* "${PEER_NODE_HOME}/peers/${PEER_HOST}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"
  verifyResult $? "Unable to copy peer tls cert."

  echo "Deleting remote node folder ..."

  ssh -i ${PEER_HOST_KEY} ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN} /bin/rm -rf ${NETWORK_HOME}/${PEER_NODE_FOLDER}

  ssh -i ${PEER_HOST_KEY} ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN} /bin/mkdir -p ${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers

  scp -C -i ${PEER_HOST_KEY} -r ${PEER_NODE_HOME}/users ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}:${NETWORK_HOME}/${PEER_NODE_FOLDER}/users
  verifyResult $? "Unable to copy peer users to ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}."

  scp -C -i ${PEER_HOST_KEY} -r ${PEER_NODE_HOME}/peers/${PEER_HOST} ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}:${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers/${PEER_HOST}
  verifyResult $? "Unable to copy peer to ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}."


  #create peer config file
	/bin/cp core-template.yaml ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml

  PEER_HOME_ESCAPED=`echo ${PEER_HOME_REMOTE} | sed -e 's/\//\\\\\//g'`
  sed -i "s/\${PEER_PORT}/${PEER_PORT}/g"                 ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml 
  verifyResult $? "Unable to update ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ."
  sed -i "s/\${PEER_MSPID}/${PEER_MSPID}/g"               ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml 
  verifyResult $? "Unable to update ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ."
  sed -i "s/\${PEER_HOME}/${PEER_HOME_ESCAPED}/g"         ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml 
  verifyResult $? "Unable to update ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ."
  sed -i "s/\${DOMAIN}/${DOMAIN}/g"                       ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml 
  verifyResult $? "Unable to update ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ."
  sed -i "s/\${PEER_HOST}/${PEER_HOST}/g"                 ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml 
  verifyResult $? "Unable to update ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ."
  sed -i "s/\${PEER_CC_PORT}/${PEER_CC_PORT}/g"           ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml 
  verifyResult $? "Unable to update ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ."
  sed -i "s/\${PEER_HOST_PREFIX}/${PEER_HOST_PREFIX}/g"   ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml 
  verifyResult $? "Unable to update ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ."

  #transfer the config file to peer machine
  scp -C -i ${PEER_HOST_KEY} ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}:${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers/${PEER_HOST}/core.yaml
  verifyResult $? "Unable to copy core.yaml to ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}."

  #transfer the external builders
  scp -C -i ${PEER_HOST_KEY} -r builders ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}:${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers/${PEER_HOST}/.
  verifyResult $? "Unable to copy external builders to ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}."

  successln "Peer Node MSP Created sussfully for peer node $1"
}

createPeerNode 1
createPeerNode 2
createPeerNode 3

