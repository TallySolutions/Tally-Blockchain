#!/bin/bash

#Usage: 6A_StartOrderServers.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 

setCANode 1

function stopOrderer()
{
   setOrdererNode $1

   PID=$(checkIfRemoteProcessRunning orderer $ORDERER_HOST_USER $ORDERER_HOST.$DOMAIN $ORDERER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      ssh -i $ORDERER_HOST_KEY $ORDERER_HOST_USER@$ORDERER_HOST.$DOMAIN "kill -9 $PID"
      PID=$(checkIfRemoteProcessRunning orderer $ORDERER_HOST_USER $ORDERER_HOST.$DOMAIN $ORDERER_HOST_KEY)
      if [[ $PID -ne -1 ]]; then
         errorln "Unable to stop orderer."
      else 
         successln "Orderer Server stopped at orderer node $1."
      fi 
   else
         warnln "Orderer server not running at orderer node $1."
   fi
}

stopOrderer 1
stopOrderer 2
