#!/bin/bash

. ./SetEnv.sh
. ./SetCANode.sh

setCANode 1

function printHelp()
{
  warnln "Usage: 19_DeleteTallyUser.sh <User Id>"          
}
if [[ $# -lt 2 ]] ; then
  printHelp
  exit 1
fi

function DeleteUser()
{
    CA_HOME=$1
    USER=$2
    export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
    infoln "Deleting User ${USER}"
           
}

DeleteUser ${TALLY_CA_HOME} $1 