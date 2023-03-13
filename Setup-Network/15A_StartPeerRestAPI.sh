#!/bin/bash

#Usage: 6A_StartOrderServers.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetPeerNode.sh 

setCANode 1

function buildRestApiService()
{
   infoln "Building Peer Rest API Service ..."
   go build -o $HOME/builds/integerkey-rest-api -C ../integerkey/application-gateway-go-rest-api 
   verifyResult $? "Build failed."
}

function startPeerRestApiService()
{
   setPeerNode $1

   PID=$(checkIfRemoteProcessRunning integerkey-rest-api $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      warnln "Peer Rest API Service already running on peer node $1 [PID: $PID], stopping first... "
      ssh -i $PEER_HOST_KEY $PEER_HOST_USER@$PEER_HOST.$DOMAIN  "kill -9 $PID"
      PID=$(checkIfRemoteProcessRunning integerkey-rest-api $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
      if [[ $PID -ne -1 ]]; then
         errorln "Unable to stop peer Rest API Service."
         return
      else 
         successln "Peer Rest API Service stopped at peer node $1."
      fi 
   fi

   infoln "Copying Peer Rest API Service binary ..."
   scp -C -i $PEER_HOST_KEY $HOME/builds/integerkey-rest-api ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}:${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers/${PEER_HOST}/.
   verifyResult $? "Copy failed."

   ssh -i $PEER_HOST_KEY $PEER_HOST_USER@$PEER_HOST.$DOMAIN "cd $PEER_HOME_REMOTE; nohup ./integerkey-rest-api > rest-api-log.txt 2>&1 &"
   PID=$(checkIfRemoteProcessRunning integerkey-rest-api $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
   if [[ $PID -ne -1 ]]; then
      successln "Peer Rest API Service started at peer node $1 [PID: $PID], you can check server log at $PEER_HOME_REMOTE/rest-api-log.txt."
   else
      errorln "Unable to start peer Rest API service."
   fi
}

buildRestApiService

startPeerRestApiService 1
startPeerRestApiService 2
startPeerRestApiService 3
