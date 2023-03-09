#!/bin/bash

#Usage: 4_CreateOrdererMSP.sh 

. ./SetGlobalVariables.sh 1

ORDERER_NODE_HOME=${TALLY_HOME}/organizations/ordererOrganizations/${DOMAIN}

mkdir -p ${ORDERER_NODE_HOME}

/bin/rm -rf ${ORDERER_NODE_HOME}/*

export FABRIC_CA_CLIENT_HOME=${ORDERER_CA_HOME}/client

#Enroll Admin
fabric-ca-client enroll -u https://${ORDERER_CA_USER}:${ORDERER_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT} --caname ${ORDERER_CA_NAME} --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=admin --tls.certfiles ${ORDERER_CA_HOME}/ca-cert.pem -M "${ORDERER_NODE_HOME}/msp"
if [[ $? -ne 0 ]]; then
	echo "Unable to enroll Orderer CA Admin MSP : is Orderer CA setup and running?"
	exit 1
fi

#Create config.yaml

DOMAIN_S=`echo ${DOMAIN} | sed -e 's/\./-/g'`
CERTNAME=${CA_HOST}-${DOMAIN_S}-${ORDERER_CA_PORT}-${ORDERER_CA_NAME}

echo "NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/${CERTNAME}.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/${CERTNAME}.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/${CERTNAME}.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/${CERTNAME}.pem
    OrganizationalUnitIdentifier: orderer" > "${ORDERER_NODE_HOME}/msp/config.yaml"

  # Copy TLS CA cert to orderer org's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${ORDERER_NODE_HOME}/msp/tlscacerts"
  cp "${TLS_CA_HOME}/ca-cert.pem" "${ORDERER_NODE_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"

  # Copy orderer org's CA cert to orderer org's /tlsca directory (for use by clients)
  mkdir -p "${ORDERER_NODE_HOME}/tlsca"
  cp "${ORDERER_CA_HOME}/ca-cert.pem" "${ORDERER_NODE_HOME}/tlsca/tlsca.${DOMAIN}-cert.pem"
:

  echo "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://${ORDERER_ADMIN_USER}:${ORDERER_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT} --caname ${ORDERER_CA_NAME} --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=admin -M "${ORDERER_NODE_HOME}/users/Admin@${DOMAIN}/msp" --tls.certfiles "${ORDERER_CA_HOME}/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${ORDERER_NODE_HOME}/msp/config.yaml" "${ORDERER_NODE_HOME}/users/Admin@${DOMAIN}/msp/config.yaml"
