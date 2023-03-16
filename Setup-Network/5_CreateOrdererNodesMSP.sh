#!/bin/bash

#Usage: 5_CreateOrdererNodeMSP.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 

setCANode 1

. ./RegisterEnroll.sh

function create()
{

  if [[ $# -lt 1 ]] ; then
    echo "Usage: create <Orderer no Number:1,2 etc.>"
    exit 1
  fi

  setOrdererNode $1

  #First Register the user
  
  RegisterUser ${ORDERER_CA_HOME} ${ORDERER_USER} ${ORDERER_PASSWORD} orderer
  
  
  ORDERER_NODE_FOLDER=organizations/ordererOrganizations/${DOMAIN}
  
  ORDERER_NODE_HOME=${TALLY_HOME}/${ORDERER_NODE_FOLDER}
  
  /bin/rm -rf "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}"
  
  export FABRIC_CA_CLIENT_HOME=${ORDERER_CA_HOME}/client

  infoln "Generating the orderer msp"

  fabric-ca-client enroll -u https://${ORDERER_USER}:${ORDERER_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT} --caname ${ORDERER_CA_NAME} --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=orderer -M "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/msp" --csr.hosts ${ORDERER_HOST}.${DOMAIN} --tls.certfiles "${ORDERER_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll orderer msp."
  
  cp "${ORDERER_NODE_HOME}/msp/config.yaml" "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  fabric-ca-client enroll -u https://${TLS_ADMIN_USER}:${TLS_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT} --caname ${TLS_CA_NAME} -M "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls" --enrollment.profile tls --csr.hosts ${ORDERER_HOST}.${DOMAIN} --tls.certfiles "${TLS_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll orderer tls."

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/tlscacerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/ca.crt"
  verifyResult $? "Unable to copy tls certificate."
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/signcerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/server.crt"
  verifyResult $? "Unable to copy tls certificate."
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/keystore/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/server.key"
  verifyResult $? "Unable to copy tls certificate."

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/msp/tlscacerts"
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/tlscacerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"
  verifyResult $? "Unable to create tlc ca certificate."


  infoln "Deleting remote node folder ..."


  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} /bin/rm -rf ${NETWORK_HOME}/${ORDERER_NODE_FOLDER}

  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} /bin/mkdir -p ${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers

  scp -C -i ${ORDERER_HOST_KEY} -r ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}:${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}
  verifyResult $? "Unable to copy ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST} to ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}"


  #create orderer config file
	/bin/cp orderer-template.yaml ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml

  ORDERER_HOME_ESCAPED=`echo ${ORDERER_HOME_REMOTE} | sed -e 's/\//\\\\\//g'`
  sed -i "s/\${ORDERER_PORT}/${ORDERER_PORT}/g"             ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_MSPID}/${ORDERER_MSPID}/g"           ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_HOME}/${ORDERER_HOME_ESCAPED}/g"     ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_ADMIN_PORT}/${ORDERER_ADMIN_PORT}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_OPS_PORT}/${ORDERER_OPS_PORT}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."


  #transfer the config file to orderer machine
  scp -C -i ${ORDERER_HOST_KEY} ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}:${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}/orderer.yaml
  verifyResult $? "Unable to copy orderer.yaml to ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}"

  successln "Successfully created orderer node $1"
}

create 1 
create 2
create 3

