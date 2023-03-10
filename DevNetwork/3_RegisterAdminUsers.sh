#!/bin/bash

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./RegisterEnroll.sh

setCANode 1

EnrollUser   ${TLS_CA_HOME} https://${TLS_CA_USER}:${TLS_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT}
RegisterUser ${TLS_CA_HOME} ${TLS_ADMIN_USER} ${TLS_ADMIN_PASSWORD} admin

EnrollUser   ${ORDERER_CA_HOME} https://${ORDERER_CA_USER}:${ORDERER_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT}
RegisterUser ${ORDERER_CA_HOME} ${ORDERER_ADMIN_USER} ${ORDERER_ADMIN_PASSWORD} admin


EnrollUser   ${TALLY_CA_HOME} https://${TALLY_CA_USER}:${TALLY_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT}
RegisterUser ${TALLY_CA_HOME} ${PEER_ADMIN_USER} ${PEER_ADMIN_PASSWORD} admin

