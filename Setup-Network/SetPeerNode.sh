#!/bin/bash

function setPeerNode()
{

  NODE=`formatNodeNo $1`

  PEER_HOST=${PEER_HOST_PREFIX}${NODE}

  PEER_HOME_REMOTE=${PEER_USER_HOME}/${NETWORK_HOME}/organizations/peerOrganizations/${DOMAIN}/peers/${PEER_HOST}

  PEER_USER=peer${NODE}
  PEER_PASSWORD=peer${NODE}pw

}