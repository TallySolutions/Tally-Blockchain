#!/bin/bash

#Usage: 7_CreatePeerMSP.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 

setCANode 1

PEER_NODE_HOME=${TALLY_HOME}/organizations/peerOrganizations/${DOMAIN}

mkdir -p ${PEER_NODE_HOME}

/bin/rm -rf ${PEER_NODE_HOME}/*

export FABRIC_CA_CLIENT_HOME=${TALLY_CA_HOME}/client

#Enroll peer
infoln "Generating the peer org msp"
fabric-ca-client enroll -u https://${TALLY_CA_USER}:${TALLY_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT} --caname ${TALLY_CA_NAME} --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=admin --tls.certfiles ${TALLY_CA_HOME}/ca-cert.pem -M "${PEER_NODE_HOME}/msp"
verifyResult $? "Unable to enroll Peer org MSP : is Tally CA setup and running?"

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
  verifyResult $? "Unable to create config.yaml."

  # Copy TLS CA cert to peer org's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PEER_NODE_HOME}/msp/tlscacerts"
  cp "${TLS_CA_HOME}/ca-cert.pem" "${PEER_NODE_HOME}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"
  verifyResult $? "Unable to create tls certificate."

  # Copy peer org's CA cert to peer org's /tlsca directory (for use by clients)
  mkdir -p "${PEER_NODE_HOME}/tlsca"
  cp "${TALLY_CA_HOME}/ca-cert.pem" "${PEER_NODE_HOME}/tlsca/tlsca.${DOMAIN}-cert.pem"
  verifyResult $? "Unable to create tls certificate."

  infoln "Generating the peer admin msp"
  fabric-ca-client enroll -u https://${PEER_ADMIN_USER}:${PEER_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT} --caname ${TALLY_CA_NAME} --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=admin -M "${PEER_NODE_HOME}/users/Admin@${DOMAIN}/msp" --tls.certfiles "${TALLY_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll Peer Admin MSP : is Tally CA setup and running?"

  cp "${PEER_NODE_HOME}/msp/config.yaml" "${PEER_NODE_HOME}/users/Admin@${DOMAIN}/msp/config.yaml"
  verifyResult $? "Unable to create config.yaml."

  successln "Peer MSP created successfully."