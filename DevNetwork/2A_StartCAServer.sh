#!/bin/bash

#Usage: 2_StartCAServer.sh 

. ./SetEnv.sh 
. ./CAUtils.sh

function startServer()
{
   CA_NAME=$1
   CA_HOME=$2

   PID=$(checkIfProcessRunning fabric-ca-server $CA_HOME)
   if [[ $PID -ne -1 ]]; then
      warnln "Fabric CA Server $CA_NAME already running [PID: $PID], skipping start."
      return 
   fi

   cd ${CA_HOME}
   fabric-ca-server start > log.txt 2>&1 &
   verifyResult $? "Could not start fabric CA Server $CA_NAME."
   successln "Fabric CA Server $CA_NAME started, you can check server log at $CA_HOME/log.txt"
}

startServer $TLS_CA_NAME $TLS_CA_HOME
startServer $TALLY_CA_NAME $TALLY_CA_HOME
startServer $ORDERER_CA_NAME $ORDERER_CA_HOME