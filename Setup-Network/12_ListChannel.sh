#!/bin/bash

#Usage: 10_ListChannel.sh 
. ./SetEnv.sh 
. ./SetOrdererNode.sh 

function list()
{
   if [[ $# -lt 1 ]] ; then
     fataln "Usage: list <Order node Number:1,2 etc.>"
   fi
   setOrdererNode $1
   
   export OSN_TLS_CA_ROOT_CERT=${TLS_CA_HOME}/ca-cert.pem
   export ADMIN_TLS_SIGN_CERT=${TALLY_HOME}/admin_client/client-tls-cert.pem
   export ADMIN_TLS_PRIVATE_KEY=${TALLY_HOME}/admin_client/client-tls-key.pem
  
   infoln "Listing channel ${CHANNEL_ID} on Orderer $1 ===>"

   osnadmin channel list --channelID ${CHANNEL_ID}  -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_ADMIN_PORT} --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY

   infoln "============== END OF OUTPUT ================="
}

list 1
list 2