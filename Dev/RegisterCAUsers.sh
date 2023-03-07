#!/bin/bash

. ./SetGlobalVariables.sh 1

function EnrollUser()
{
        CA_HOME=$1
        URL=$2
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        echo "Enrolling user for ${CA_HOME} [${TYPE}]"
        set -x
        fabric-ca-client enroll -u ${URL} --tls.certfiles "${CA_HOME}/ca-cert.pem"
        { set +x; } 2>/dev/null

}
function RegisterUser()
{
        CA_HOME=$1
        USER=$2
        PASSWORD=$3
        TYPE=$4
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        echo "Registering user ${USER} for ${CA_HOME} [${TYPE}]"
        set -x
        fabric-ca-client register --id.name ${USER} --id.secret ${PASSWORD} --id.type ${TYPE} --tls.certfiles "${CA_HOME}/ca-cert.pem"
        { set +x; } 2>/dev/null

}

EnrollUser   ${TLS_CA_HOME} https://${TLS_CA_USER}:${TLS_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT}
RegisterUser ${TLS_CA_HOME} ${TLS_ADMIN_USER} ${TLS_ADMIN_PASSWORD} admin

EnrollUser   ${ORDERER_CA_HOME} https://${ORDERER_CA_USER}:${ORDERER_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT}
RegisterUser ${ORDERER_CA_HOME} ${ORDERER_USER} ${ORDERER_PASSWORD} orderer
RegisterUser ${ORDERER_CA_HOME} ${ORDERER_ADMIN_USER} ${ORDERER_ADMIN_PASSWORD} admin


EnrollUser   ${TALLY_CA_HOME} https://${TALLY_CA_USER}:${TALLY_CA_PASSWORD}@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT}
RegisterUser ${TALLY_CA_HOME} ${PEER_USER} ${PEER_PASSWORD} peer
RegisterUser ${TALLY_CA_HOME} ${PEER_ADMIN_USER} ${PEER_ADMIN_PASSWORD} admin

