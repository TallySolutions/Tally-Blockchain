#!/bin/bash

. ./SetEnv.sh
. ./SetCANode.sh

setCANode 1

function printHelp()
{
  warnln "Usage: 17_CreateTallyClientUser.sh <User Id> <Is_Approver> <Is_Creator>"               # How do we get CA_HONE and TYPE values? Ans: line 31 are params to function, the input params in line 31= those from cli
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
        APPROVER=$4
        CREATOR=$5
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        infoln "Registering user ${USER} for ${CA_HOME} [${TYPE}] \n Is approver?= ${APPROVER} \n Is creator?= ${CREATOR}"
        fabric-ca-client register --id.name ${USER} --id.type ${TYPE} --id.affiliation tally --id.attrs 'approver='$APPROVER':ecert, creator:'$CREATOR':ecert' --tls.certfiles "${CA_HOME}/ca-cert.pem"  # added attributes for id affiliations and other attrs
        res=$?
        verifyResult $res "Unable to register user for ${CA_HOME} [${TYPE}]. (Are you trying to re-register a user or CA server is stopped?)"

}

RegisterUser ${TALLY_CA_HOME} $1 client $2 $3


# creator and approver- categories for users