#!/bin/bash

#Usage: 5_CreateOrdererNodeMSP.sh <ORDERER-NODE-NO>

if [[ $# -lt 1 ]] ; then
  echo "Usage: 5_CreateOrdererNodeMSP.sh <Orderer no Number:1,2 etc.>"
  exit 1
fi

. ./SetGlobalVariables.sh $1

ORDERER_NODE_FOLDER=organizations/ordererOrganizations/${DOMAIN}

ORDERER_NODE_HOME=${TALLY_HOME}/${ORDERER_NODE_FOLDER}

/bin/rm -rf "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}"

export FABRIC_CA_CLIENT_HOME=${ORDERER_CA_HOME}/client

  echo "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://${ORDERER_USER}:${ORDERER_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT} --caname ${ORDERER_CA_NAME} -M "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/msp" --csr.hosts ${ORDERER_HOST}.${DOMAIN} --tls.certfiles "${ORDERER_CA_HOME}/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${ORDERER_NODE_HOME}/msp/config.yaml" "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/msp/config.yaml"

  echo "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://${TLS_ADMIN_USER}:${TLS_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT} --caname ${TLS_CA_NAME} -M "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls" --enrollment.profile tls --csr.hosts ${ORDERER_HOST}.${DOMAIN} --tls.certfiles "${TLS_CA_HOME}/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/tlscacerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/ca.crt"
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/signcerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/server.crt"
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/keystore/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/server.key"

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/msp/tlscacerts"
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/tls/tlscacerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"



  echo "Deleting remote node folder ..."

set -x

  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} /bin/rm -rf ${NETWORK_HOME}/${ORDERER_NODE_FOLDER}

  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} /bin/mkdir -p ${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers

  scp -C -i ${ORDERER_HOST_KEY} -r ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}:${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}

{ set +x; } 2>/dev/null


  #create orderer config file
	/bin/cp orderer-template.yaml ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml

  ORDERER_HOME_ESCAPED=`echo ${ORDERER_HOME_REMOTE} | sed -e 's/\//\\\\\//g'`
  sed -e "s/\${ORDERER_PORT}/${ORDERER_PORT}/g"             ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml   > ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.1 
  sed -e "s/\${ORDERER_MSPID}/${ORDERER_MSPID}/g"           ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.1 > ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.2 
  sed -e "s/\${ORDERER_HOME}/${ORDERER_HOME_ESCAPED}/g"     ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.2 > ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.3 
  sed -e "s/\${ORDERER_ADMIN_PORT}/${ORDERER_ADMIN_PORT}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.3 > ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.4 

  /bin/rm ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.1 ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.2 ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.3
  /bin/mv ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml.4  ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml

  #transfer the config file to orderer machine
  scp -C -i ${ORDERER_HOST_KEY} ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}/orderer.yaml ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}:${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}/orderer.yaml



