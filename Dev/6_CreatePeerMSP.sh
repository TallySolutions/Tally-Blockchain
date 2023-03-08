#!/bin/bash

#Usage: 6_CreatePeerMSP.sh 

. ./SetGlobalVariables.sh 1

PEER_NODE_HOME=${TALLY_HOME}/organizations/peerOrganizations/${DOMAIN}

mkdir -p ${PEER_NODE_HOME}

/bin/rm -rf ${PEER_NODE_HOME}/*

export FABRIC_CA_CLIENT_HOME=${TALLY_CA_HOME}/client

#Enroll Admin
fabric-ca-client enroll -u https://${TALLY_CA_USER}:${TALLY_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT} --caname ${TALLY_CA_NAME} --tls.certfiles ${TALLY_CA_HOME}/ca-cert.pem -M "${PEER_NODE_HOME}/msp"
if [[ $? -ne 0 ]]; then
	echo "Unable to enroll Peer CA Admin MSP : is Tally CA setup and running?"
	exit 1
fi

#Create config.yaml

DOMAIN_S=`echo ${DOMAIN} | sed -e 's/\./-/g'`
CERTNAME=${CA_HOST}-${DOMAIN_S}-${TALLY_CA_PORT}-${TALLY_CA_NAME}

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
    OrganizationalUnitIdentifier: orderer" > "${PEER_NODE_HOME}/msp/config.yaml"

  # Copy TLS CA cert to peer org's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PEER_NODE_HOME}/msp/tlscacerts"
  cp "${TLS_CA_HOME}/ca-cert.pem" "${PEER_NODE_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"

  # Copy peer org's CA cert to peer org's /tlsca directory (for use by clients)
  mkdir -p "${PEER_NODE_HOME}/tlsca"
  cp "${TALLY_CA_HOME}/ca-cert.pem" "${PEER_NODE_HOME}/tlsca/tlsca.${DOMAIN}-cert.pem"
:

  echo "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://${PEER_ADMIN_USER}:${PEER_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT} --caname ${TALLY_CA_NAME} -M "${PEER_NODE_HOME}/users/Admin@${DOMAIN}/msp" --tls.certfiles "${TALLY_CA_HOME}/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PEER_NODE_HOME}/msp/config.yaml" "${PEER_NODE_HOME}/users/Admin@${DOMAIN}/msp/config.yaml"
