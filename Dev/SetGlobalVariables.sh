#!/bin/bash


if [[ $# -lt 1 ]] ; then
  echo "Usage: SetGlobalVariables.sh <node_number>"
fi

NODE=$1

re='^[0-9]+$'
if ! [[ $NODE =~ $re ]] ; then
   NODE=1
fi

if [[ ${#NODE} -lt 10 ]] ; then
    NODE="0${#NODE}"
fi


TALLY_HOME=/home/ubuntu/fabric/tally-network

DOMAIN=tally.tallysolutions.com

#Hosts

CA_HOST=tbchlfdevca01

ORDERER_HOST=tbchlfdevord${NODE}

PEER_HOST=tbchlfdevpeer${NODE}

#CA Servers

CA_SERVER_HOME=${TALLY_HOME}/fabric-ca-servers

#TLS CA Server

TLS_CA_NAME=tls
TLS_CA_HOME=${CA_SERVER_HOME}/${TLS_CA_NAME}
TLS_CA_PORT=7054
TLS_CA_OPS_PORT=9443
TLS_CA_USER=tlsadmin
TLS_CA_PASSWORD=tlsadminpw

#Tally CA 
TALLY_CA_NAME=tally
TALLY_CA_HOME=${CA_SERVER_HOME}/${TALLY_CA_NAME}
TALLY_CA_PORT=7055
TALLY_CA_OPS_PORT=9444
TALLY_CA_USER=tallyadmin
TALLY_CA_PASSWORD=tallyadminpw

#Orderer CA 
ORDERER_CA_NAME=orderer
ORDERER_CA_HOME=${CA_SERVER_HOME}/${ORDERER_CA_NAME}
ORDERER_CA_PORT=7056
ORDERER_CA_OPS_PORT=9445
ORDERER_CA_USER=ordadmin
ORDERER_CA_PASSWORD=ordadminpw

#Organizations

ORG_HOME=${TALLY_HOME}/organizations
ORDERER_HOME=${ORG_HOME}/ordererOrganizations/${DOMAIN}
PEER_HOME=${ORG_HOME}/ordererOrganizations/${DOMAIN}

ORDERER_NODE_HOME=${ORDERER_HOME}/orderers/${ORDERER_HOST}

PEER_NODE_HOME=${PEER_HOME}/peers/${PEER_HOST}


#USERS

TLS_ADMIN_USER=rcaadmin
TLS_ADMIN_PASSWORD=rcadminpw

ORDERER_USER=orderer
ORDERER_PASSWORD=ordererpw

ORDERER_ADMIN_USER=ordererAdmin
ORDERER_ADMIN_PASSWORD=ordererAdminpw


PEER_USER=peer
PEER_PASSWORD=peerpw

PEER_ADMIN_USER=peerAdmin
PEER_ADMIN_PASSWORD=peerAdminpw

