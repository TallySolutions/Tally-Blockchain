#!/bin/bash

#Usage: 6A_StartOrderServers.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetPeerNode.sh 

setCANode 1

function stopPeer()
{
   setPeerNode $1

   PID=$(checkIfRemoteProcessRunning peer $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      ssh -i $PEER_HOST_KEY $PEER_HOST_USER@$PEER_HOST.$DOMAIN  "kill -9 $PID"
      PID=$(checkIfRemoteProcessRunning peer $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
      if [[ $PID -ne -1 ]]; then
         errorln "Unable to stop peer."
      else 
         successln "Peer Server stopped at peer node $1."
      fi 
   else
         warnln "Peer server not running at peer node $1."
   fi
}

stopPeer 1
stopPeer 2
stopPeer 3
