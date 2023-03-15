#!/bin/bash



function setOrdererNode()
{

  NODE=`formatNodeNo $1`

  ORDERER_HOST=${ORDERER_HOST_PREFIX}${NODE}
  ORDERER_HOME_REMOTE=${ORDERER_USER_HOME}/${NETWORK_HOME}/organizations/ordererOrganizations/${DOMAIN}/orderers/${ORDERER_HOST}

  ORDERER_USER=orderer${NODE}
  ORDERER_PASSWORD=orderer${NODE}pw


}


function setOrdererDockerNode()
{

  NODE=`formatNodeNo $1`

  ORDERER_HOST=${ORDERER_HOST_PREFIX}${NODE}
  ORDERER_HOME_REMOTE=${ORDERER_USER_HOME}/${NETWORK_HOME}/organizations/ordererOrganizations/${DOMAIN}/orderers/${ORDERER_HOST}_$2

  ORDERER_USER=orderer${NODE}P$2
  ORDERER_PASSWORD=orderer${NODE}P${2}pw


}

export -f setOrdererNode
export -f setOrdererDockerNode