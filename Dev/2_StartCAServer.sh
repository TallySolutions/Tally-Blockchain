#!/bin/bash

#Usage: 2_StartCAServer.sh 

#pass default node as 1, node is not used in this script
. ./SetGlobalVariables.sh 1

function start()
{

    if [[ $# -lt 1 ]] ; then
      echo "Usage: start <TLS|TALLY|ORDERER>"
      exit 1
    fi
    
    if [[ $1 == "TLS" ]]; then
       cd ${TLS_CA_HOME}
       fabric-ca-server start > log.txt 2>&1 &
    fi
    
    if [[ $1 == "TALLY" ]]; then
       cd ${TALLY_CA_HOME}
       fabric-ca-server start &> log.txt 2>&1 &
    fi
    
    if [[ $1 == "ORDERER" ]]; then
       cd ${ORDERER_CA_HOME}
       fabric-ca-server start &> log.txt 2>&1 &
    fi
}

start TLS
start ORDERER
start TALLY