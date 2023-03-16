#!/bin/bash

function setPeerNode()
{

  NODE=`formatNodeNo $1`

  PEER_HOST=${PEER_HOST_PREFIX}${NODE}

  PEER_HOME_REMOTE=${PEER_USER_HOME}/${NETWORK_HOME}/organizations/peerOrganizations/${DOMAIN}/peers/${PEER_HOST}

  PEER_USER=peer${NODE}
  PEER_PASSWORD=peer${NODE}pw

}

function setup_peer_paths()
{
  export FABRIC_CFG_PATH=${PEER_HOME}/peers/${PEER_HOST}
  export CORE_PEER_MSPCONFIGPATH=${PEER_HOME}/users/Admin@${DOMAIN}/msp
  export CORE_PEER_ADDRESS=${PEER_HOST}.${DOMAIN}:${PEER_PORT}
  export CORE_PEER_LOCALMSPID=${PEER_MSPID}
  export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_HOME}/peers/${PEER_HOST}/tls/ca.crt
}
