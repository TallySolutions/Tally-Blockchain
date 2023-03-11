#!/bin/bash

#Usage: 2B_StopCAServer.sh 

. ./SetEnv.sh 
. ./CAUtils.sh

function stopServer()
{
   CA_NAME=$1
   CA_HOME=$2

   PID=$(checkIfProcessRunning fabric-ca-server $CA_HOME)
   if [[ $PID -ne -1 ]]; then
      infoln "Stopping server $CA_NAME [PID: $PID] ..."
      kill -9 $PID
      verifyResult $? "Could not kill PID $PID."
      successln "Fabric CA Server $CA_NAME stopped."
      return 
   fi

   warnln "Fabric CA Server $CA_NAME not running."
}

stopServer $TLS_CA_NAME $TLS_CA_HOME
stopServer $TALLY_CA_NAME $TALLY_CA_HOME
stopServer $ORDERER_CA_NAME $ORDERER_CA_HOME