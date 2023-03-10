#!/bin/bash

. ./Utils.sh

function setEnv()
{

 #Global

 NETWORK_HOME=fabric/tally-network
 TALLY_HOME=$HOME/$NETWORK_HOME
 DOMAIN=tally.tallysolutions.com
 
 #Hosts

 CA_HOST_PREFIX=tbchlfdevca
 ORDERER_HOST_PREFIX=tbchlfdevord
 PEER_HOST_PREFIX=tbchlfdevpeer
 
 ORDERER_HOST_KEY=$HOME/.ssh/TDevBC-Orderer-1-KeyPair.pem
 PEER_HOST_KEY=$HOME/.ssh/TDevBC-Peer-keyPair.pem
 
 ORDERER_HOST_USER=ubuntu
 PEER_HOST_USER=ubuntu
 
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
 

 #CA USERS
 
 TLS_ADMIN_USER=rcaadmin
 TLS_ADMIN_PASSWORD=rcadminpw
 

 ORDERER_ADMIN_USER=ordererAdmin
 ORDERER_ADMIN_PASSWORD=ordererAdminpw
 
 
 PEER_ADMIN_USER=peerAdmin
 PEER_ADMIN_PASSWORD=peerAdminpw
 
 #Organizations
 
 ORG_HOME=${TALLY_HOME}/organizations
 ORDERER_HOME=${ORG_HOME}/ordererOrganizations/${DOMAIN}
 PEER_HOME=${ORG_HOME}/peerOrganizations/${DOMAIN}


 #ORDERER
 
 ORDERER_PORT=7060
 ORDERER_ADMIN_PORT=9443
 ORDERER_MSPID=Orderer
 ORDERER_USER_HOME=/home/ubuntu
 
 PEER_PORT=7051
 PEER_CC_PORT=7052
 PEER_USER_HOME=/home/ubuntu
 PEER_MSPID=Tally
 
 #Channel
 
 CHANNEL_ID=integerkey
 
 #Chaincode
 
 CC_RUNTIME_LANGUAGE=golang

}

export -f setEnv

. ./SetOrdererNode.sh
. ./SetPeerNode.sh


infoln "Setting up environment ..."
setEnv