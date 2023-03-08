#!/bin/bash

#Usage: 10_SetupChannels.sh <ORDERER-NODE-NO>

if [[ $# -lt 1 ]] ; then
  echo "Usage: 10_SetupChannels.sh <Order node Number:1,2 etc.>"
  exit 1
fi
. ./SetGlobalVariables.sh $1

export OSN_TLS_CA_ROOT_CERT=${TLS_CA_HOME}/ca-cert.pem
export ADMIN_TLS_SIGN_CERT=${TALLY_HOME}/admin_client/client-tls-cert.pem
export ADMIN_TLS_PRIVATE_KEY=${TALLY_HOME}/admin_client/client-tls-key.pem


 #create configtx file
 /bin/cp configtx-template.yaml ${TALLY_HOME}/admin_client/configtx.yaml

  ORDERER_HOME_ESCAPED=`echo ${ORDERER_HOME} | sed -e 's/\//\\\\\//g'`
  PEER_HOME_ESCAPED=`echo ${PEER_HOME} | sed -e 's/\//\\\\\//g'`
  sed -e "s/\${ORDERER_MSPID}/${ORDERER_MSPID}/g"             ${TALLY_HOME}/admin_client/configtx.yaml   > ${TALLY_HOME}/admin_client/configtx.yaml.1 
  sed -e "s/\${ORDERER_HOME}/${ORDERER_HOME_ESCAPED}/g"       ${TALLY_HOME}/admin_client/configtx.yaml.1 > ${TALLY_HOME}/admin_client/configtx.yaml.2 
  sed -e "s/\${ORDERER_HOST_PREFIX}/${ORDERER_HOST_PREFIX}/g" ${TALLY_HOME}/admin_client/configtx.yaml.2 > ${TALLY_HOME}/admin_client/configtx.yaml.3 
  sed -e "s/\${DOMAIN}/${DOMAIN}/g"                           ${TALLY_HOME}/admin_client/configtx.yaml.3 > ${TALLY_HOME}/admin_client/configtx.yaml.4 
  sed -e "s/\${ORDERER_PORT}/${ORDERER_PORT}/g"               ${TALLY_HOME}/admin_client/configtx.yaml.4 > ${TALLY_HOME}/admin_client/configtx.yaml.5 
  sed -e "s/\${PEER_MSPID}/${PEER_MSPID}/g"                   ${TALLY_HOME}/admin_client/configtx.yaml.5 > ${TALLY_HOME}/admin_client/configtx.yaml.6 
  sed -e "s/\${PEER_HOME}/${PEER_HOME_ESCAPED}/g"             ${TALLY_HOME}/admin_client/configtx.yaml.6 > ${TALLY_HOME}/admin_client/configtx.yaml.7 

  /bin/rm ${TALLY_HOME}/admin_client/configtx.yaml ${TALLY_HOME}/admin_client/configtx.yaml.1 ${TALLY_HOME}/admin_client/configtx.yaml.2 ${TALLY_HOME}/admin_client/configtx.yaml.3 ${TALLY_HOME}/admin_client/configtx.yaml.4 ${TALLY_HOME}/admin_client/configtx.yaml.5 ${TALLY_HOME}/admin_client/configtx.yaml.6
  /bin/mv ${TALLY_HOME}/admin_client/configtx.yaml.7  ${TALLY_HOME}/admin_client/configtx.yaml

export FABRIC_CFG_PATH=${TALLY_HOME}/admin_client/

mkdir -p FABRIC_CFG_PATH=${TALLY_HOME}/admin_client/${CHANNEL_ID}

set -x

configtxgen -profile TallyApplicationGenesis -outputBlock ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb -channelID ${CHANNEL_ID}

osnadmin channel join --channelID ${CHANNEL_ID}  --config-block ${TALLY_HOME}/admin_client/${CHANNEL_ID}/genesis_block.pb -o ${ORDERER_HOST}.${DOMAIN}:${ORDERER_ADMIN_PORT} --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY

set +x 2> /dev/null