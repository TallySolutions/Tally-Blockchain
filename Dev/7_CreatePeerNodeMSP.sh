#!/bin/bash

#Usage: 7_CreatePeerNodeMSP.sh <PEER-NODE-NO>

if [[ $# -lt 1 ]] ; then
  echo "Usage: 7_CreatePeerNodeMSP.sh <Peer no Number:1,2 etc.>"
  exit 1
fi

. ./SetGlobalVariables.sh $1

PEER_NODE_FOLDER=organizations/peerOrganizations/${DOMAIN}

PEER_NODE_HOME=${TALLY_HOME}/${PEER_NODE_FOLDER}

/bin/rm -rf "${PEER_NODE_HOME}/peers/${PEER_HOST}"

export FABRIC_CA_CLIENT_HOME=${TALLY_CA_HOME}/client

  echo "Generating the peer msp"
  set -x
  fabric-ca-client enroll -u https://${PEER_USER}:${PEER_PASSWORD}@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT} --caname ${TALLY_CA_NAME} -M "${PEER_NODE_HOME}/peers/${PEER_HOST}/msp" --csr.hosts ${PEER_HOST}.${DOMAIN} --tls.certfiles "${TALLY_CA_HOME}/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PEER_NODE_HOME}/msp/config.yaml" "${PEER_NODE_HOME}/peers/${PEER_HOST}/msp/config.yaml"

  echo "Generating the peer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://${TLS_ADMIN_USER}:${TLS_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT} --caname ${TLS_CA_NAME} -M "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls" --enrollment.profile tls --csr.hosts ${PEER_HOST}.${DOMAIN} --tls.certfiles "${TLS_CA_HOME}/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/tlscacerts/"* "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/ca.crt"
  cp "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/signcerts/"* "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/server.crt"
  cp "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/keystore/"* "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/server.key"

  # Copy peer org's CA cert to peer's /msp/tlscacerts directory (for use in the peer MSP definition)
  mkdir -p "${PEER_NODE_HOME}/peers/${PEER_HOST}/msp/tlscacerts"
  cp "${PEER_NODE_HOME}/peers/${PEER_HOST}/tls/tlscacerts/"* "${PEER_NODE_HOME}/peers/${PEER_HOST}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"



  echo "Deleting remote node folder ..."

set -x

  ssh -i ${PEER_HOST_KEY} ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN} /bin/rm -rf ${NETWORK_HOME}/${PEER_NODE_FOLDER}

  ssh -i ${PEER_HOST_KEY} ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN} /bin/mkdir -p ${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers

  scp -C -i ${PEER_HOST_KEY} -r ${PEER_NODE_HOME}/peers/${PEER_HOST} ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}:${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers/${PEER_HOST}

{ set +x; } 2>/dev/null


  #create peer config file
	/bin/cp core-template.yaml ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml

  PEER_HOME_ESCAPED=`echo ${PEER_HOME_REMOTE} | sed -e 's/\//\\\\\//g'`
  sed -e "s/\${PEER_PORT}/${PEER_PORT}/g"         ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml   > ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.1 
  sed -e "s/\${PEER_MSPID}/${PEER_MSPID}/g"       ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.1 > ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.2 
  sed -e "s/\${PEER_HOME}/${PEER_HOME_ESCAPED}/g" ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.2 > ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.3 
  sed -e "s/\${DOMAIN}/${DOMAIN}/g"               ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.3 > ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.4 
  sed -e "s/\${PEER_HOST}/${PEER_HOST}/g"         ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.4 > ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.5 
  sed -e "s/\${PEER_CC_PORT}/${PEER_CC_PORT}/g"   ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.5 > ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.6 
  sed -e "s/\${PEER_HOST_PREFIX}/${PEER_HOST_PREFIX}/g"   ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.6 > ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.7 

  /bin/rm ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.1 ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.2 ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.3 ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.4 ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.5 ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.6
  /bin/mv ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml.7  ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml

  #transfer the config file to peer machine
  scp -C -i ${PEER_HOST_KEY} ${PEER_NODE_HOME}/peers/${PEER_HOST}/core.yaml ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}:${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers/${PEER_HOST}/core.yaml



