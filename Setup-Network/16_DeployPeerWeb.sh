#!/bin/bash

#Usage: 6A_StartOrderServers.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetPeerNode.sh 

setCANode 1

function packagePeerWebServer()
{
   infoln "Packaging Peer Web Server ..."
   mkdir -p $HOME/builds
   tar -czf $HOME/builds/integerkeyapp-react-web.tar.gz --exclude="node_modules" ../integerkey/integerkeyapp-react-web 2> /dev/null
   verifyResult $? "Packaging failed."
}

function startPeerWebServer()
{
   setPeerNode $1

   infoln "Copying Peer Web Server Package ..."
   ssh -i $PEER_HOST_KEY ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN} "mkdir -p ${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers/${PEER_HOST}/web"
   scp -C -i $PEER_HOST_KEY $HOME/builds/integerkeyapp-react-web.tar.gz ${PEER_HOST_USER}@${PEER_HOST}.${DOMAIN}:${NETWORK_HOME}/${PEER_NODE_FOLDER}/peers/${PEER_HOST}/web/.
   verifyResult $? "Copy failed."

   infoln "Extracting Peer Web Server Package ..."
   ssh -i $PEER_HOST_KEY $PEER_HOST_USER@$PEER_HOST.$DOMAIN "cd $PEER_HOME_REMOTE/web; tar -xzf integerkeyapp-react-web.tar.gz"
   verifyResult $? "Extraction failed."

   PID=$(checkIfRemoteProcessRunning npm $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
   if [[ $PID -eq -1 ]]; then
      infoln "Starting Peer Web Service on peer node $1 ..."
      ssh -i $PEER_HOST_KEY $PEER_HOST_USER@$PEER_HOST.$DOMAIN "cd $PEER_HOME_REMOTE/web/integerkey/integerkeyapp-react-web; npm list ract || npm install react > $PEER_HOME_REMOTE/web/log.txt 2>&1; npm list react-id-generator || npm install react-id-generator  >> $PEER_HOME_REMOTE/web/log.txt 2>&1; npm list react-icons || npm install react-icons >> $PEER_HOME_REMOTE/web/log.txt 2>&1; nohup npm start >> $PEER_HOME_REMOTE/web/log.txt 2>&1 &"
      PID=$(checkIfRemoteProcessRunning npm $PEER_HOST_USER $PEER_HOST.$DOMAIN $PEER_HOST_KEY)
      if [[ $PID -ne -1 ]]; then
          successln "Peer Web Server started at peer node $1 [PID: $PID], you can check server log at $PEER_HOME_REMOTE/web/log.txt."
      else
          errorln "Unable to start peer Web Server."
      fi
   else
      successln "Peer Web content updated at peer node $1 - web server already running [PID: $PID], you can check server log at $PEER_HOME_REMOTE/web/log.txt."
   fi
}

packagePeerWebServer

startPeerWebServer 1
startPeerWebServer 2
startPeerWebServer 3
