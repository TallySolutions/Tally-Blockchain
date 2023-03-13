#!/bin/bash

#Usage: 6A_StartOrderServers.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetPeerNode.sh 

setCANode 1

function stopPeerRestAPISvc()
{
   setPeerNode $1

   PID=$(checkIfRemoteProcessRunning integerkey-rest-api $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      ssh -i $PEER_HOST_KEY $PEER_HOST_USER@$PEER_HOST.$DOMAIN  "kill -9 $PID"
      PID=$(checkIfRemoteProcessRunning integerkey-rest-api $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
      if [[ $PID -ne -1 ]]; then
         errorln "Unable to stop peer Rest API Service."
      else 
         successln "Peer Rest API Service stopped at peer node $1."
      fi 
   else
         warnln "Peer Rest API Service not running at peer node $1."
   fi
}

stopPeerRestAPISvc 1
stopPeerRestAPISvc 2
stopPeerRestAPISvc 3
