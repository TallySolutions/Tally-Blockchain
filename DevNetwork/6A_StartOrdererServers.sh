#!/bin/bash

#Usage: 6A_StartOrderServers.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 

setCANode 1

function startOrderer()
{
   setOrdererNode $1

   PID=$(checkIfRemoteProcessRunning orderer $ORDERER_HOST_USER $ORDERER_HOST.$DOMAIN $ORDERER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      warnln "Orderer server already running on orderer node $1 [PID: $PID], skipping start."
      return 
   fi

   ssh -i $ORDERER_HOST_KEY $ORDERER_HOST_USER@$ORDERER_HOST.$DOMAIN ". .profile; cd $ORDERER_HOME_REMOTE ; nohup orderer > log.txt 2>&1 &"
   PID=$(checkIfRemoteProcessRunning orderer $ORDERER_HOST_USER $ORDERER_HOST.$DOMAIN $ORDERER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      successln "Orderer Server started at orderer node $1 [PID: $PID], you can check server log at $ORDERER_HOME_REMOTE/log.txt."
   else
      errorln "Unable to start orderer."
   fi
}

startOrderer 1
startOrderer 2
