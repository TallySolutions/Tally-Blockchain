#!/bin/bash

. ./SetEnv.sh
. ./SetCANode.sh

setCANode 1

function printHelp()
{
  warnln "Usage: 18_EnrollTallyUser.sh <User Id> <Password>"
}
if [[ $# -lt 2 ]] ; then
  printHelp
  exit 1
fi

function EnrollUser()
{
        CA_HOME=$1
        URL=$2
	      MSP=$3
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        infoln "Enrolling user for ${CA_HOME}"
        fabric-ca-client enroll -u ${URL}  --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=client -M "${MSP}" --tls.certfiles "${CA_HOME}/ca-cert.pem"
        verifyResult $? "Unable to enroll user for ${CA_HOME}. (Is relevant CA server running?)"

}


EnrollUser ${TALLY_CA_HOME} https://$1:$2@${CA_HOST}.${DOMAIN}:${TALLY_CA_PORT} "${CLIENT_HOME}/users/$1/msp"
