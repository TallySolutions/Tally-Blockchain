#!/bin/bash

. ./SetEnv.sh
. ./SetCANode.sh

setCANode 1

function printHelp()
{
  warnln "Usage: 17_CreateTallyClientUser.sh <User Id>"
}
if [[ $# -lt 1 ]] ; then
  printHelp
  exit 1
fi

function RegisterUser()
{
        CA_HOME=$1
        USER=$2
        TYPE=$3
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        infoln "Registering user ${USER} for ${CA_HOME} [${TYPE}]"
        fabric-ca-client register --id.name ${USER} --id.type ${TYPE} --tls.certfiles "${CA_HOME}/ca-cert.pem"
        res=$?
        verifyResult $res "Unable to register user for ${CA_HOME} [${TYPE}]. (Are you trying to re-register a user or CA server is stopped?)"

}

RegisterUser ${TALLY_CA_HOME} $1 client