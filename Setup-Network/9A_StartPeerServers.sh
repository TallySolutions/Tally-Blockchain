#!/bin/bash

#Usage: 6A_StartOrderServers.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetPeerNode.sh 

setCANode 1

function startPeer()
{
   setPeerNode $1

   PID=$(checkIfRemoteProcessRunning peer $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      warnln "Peer server already running on orderer node $1 [PID: $PID], skipping start."
      return 
   fi

   ssh -i $PEER_HOST_KEY $PEER_HOST_USER@$PEER_HOST.$DOMAIN ". .profile; cd $PEER_HOME_REMOTE ; nohup peer node start > log.txt 2>&1 &"
   PID=$(checkIfRemoteProcessRunning peer $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      successln "Peer Server started at orderer node $1 [PID: $PID], you can check server log at $PEER_HOME_REMOTE/log.txt."
   else
      errorln "Unable to start peer."
   fi
}

startPeer 1
startPeer 2
startPeer 3
