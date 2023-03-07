#!/bin/bash

#Usage: CreateOrdererNodeMSP.sh <ORDERER-NODE-NO>

if [[ $# -lt 1 ]] ; then
  echo "Usage: CreateOrdererNodeMSP.sh <Orderer no Number:1,2 etc.>"
fi

. ./SetGlobalVariables.sh $1

ORDERER_NODE_HOME=${TALLY_HOME}/organizations/ordererOrganizations/${DOMAIN}

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

  echo "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://${ORDERER_ADMIN_USER}:${ORDERER_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT} --caname ${ORDERER_CA_NAME} -M "${ORDERER_NODE_HOME}/users/Admin@${DOMAIN}/msp" --tls.certfiles "${ORDERER_CA_HOME}/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${ORDERER_NODE_HOME}/msp/config.yaml" "${ORDERER_NODE_HOME}/users//Admin@${DOMAIN}/msp/config.yaml"

