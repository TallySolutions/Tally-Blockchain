#!/bin/bash

#Usage: 2_StartCAServer.sh <CA Server Name> 

#pass default node as 1, node is not used in this script
. ./SetGlobalVariables.sh 1

if [[ $# -lt 1 ]] ; then
  echo "Usage: 2_StartCAServer.sh <TLS|TALLY|ORDERER>"
  exit 1
fi

if [[ $1 == "TLS"]]; then
   cd ${TLS_CA_HOME}
   fabric-ca-server start
fi

if [[ $1 == "TALLY"]]; then
   cd ${TALLY_CA_HOME}
   fabric-ca-server start
fi

if [[ $1 == "ORDERER"]]; then
   cd ${ORDERER_CA_HOME}
   fabric-ca-server start
fi
