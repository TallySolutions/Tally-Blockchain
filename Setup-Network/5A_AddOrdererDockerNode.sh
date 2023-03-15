#!/bin/bash

#Usage: 5A_AddOrdererDockerNode.sh <physical node no> <orderer_port> 


. ./SetEnv.sh 
. ./SetCANode.sh 
. ./SetOrdererNode.sh 

  if [[ $# -lt 4 ]] ; then
    fatalln "Usage: 5A_AddOrdererDockerNode.sh <Orderer node Number:1,2 etc.> <orderer_port> <orderer_admin_port> <orderer_operation_port> [--no-reg]"
  fi

  PHYSICAL_SVR_NODE=$1
  ORDERER_PORT=$2
  ORDERER_ADMIN_PORT=$3
  ORDERER_OPS_PORT=$4
  NO_REG=$5

setCANode 1

. ./RegisterEnroll.sh

function createDockerCompose()
{

  #First Register the user
  if [[ "$NO_REG" == "--no-reg" ]]; then
     warnln "Skipping registering user ${ORDERER_USER}."
  else 
     RegisterUser ${ORDERER_CA_HOME} ${ORDERER_USER} ${ORDERER_PASSWORD} orderer
  fi
  
  ORDERER_NODE_FOLDER=organizations/ordererOrganizations/${DOMAIN}
  
  ORDERER_NODE_HOME=${TALLY_HOME}/${ORDERER_NODE_FOLDER}
  
  /bin/rm -rf "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}"
  
  export FABRIC_CA_CLIENT_HOME=${ORDERER_CA_HOME}/client

  infoln "Generating the orderer msp"

  fabric-ca-client enroll -u https://${ORDERER_USER}:${ORDERER_PASSWORD}@${CA_HOST}.${DOMAIN}:${ORDERER_CA_PORT} --caname ${ORDERER_CA_NAME} --csr.names C=IN,ST=Bengaluru,L=Bengaluru,O=Tally,OU=orderer -M "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/msp" --csr.hosts ${ORDERER_HOST}.${DOMAIN} --tls.certfiles "${ORDERER_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll orderer msp."
  
  cp "${ORDERER_NODE_HOME}/msp/config.yaml" "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  fabric-ca-client enroll -u https://${TLS_ADMIN_USER}:${TLS_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT} --caname ${TLS_CA_NAME} -M "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/tls" --enrollment.profile tls --csr.hosts ${ORDERER_HOST}.${DOMAIN} --tls.certfiles "${TLS_CA_HOME}/ca-cert.pem"
  verifyResult $? "Unable to enroll orderer tls."

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/tls/tlscacerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/tls/ca.crt"
  verifyResult $? "Unable to copy tls certificate."
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/tls/signcerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/tls/server.crt"
  verifyResult $? "Unable to copy tls certificate."
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/tls/keystore/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/tls/server.key"
  verifyResult $? "Unable to copy tls certificate."

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/msp/tlscacerts"
  cp "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/tls/tlscacerts/"* "${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"
  verifyResult $? "Unable to create tlc ca certificate."


  infoln "Deleting remote node folder ..."


  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} /bin/rm -rf ${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}_${ORDERER_PORT}

  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} /bin/mkdir -p ${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers

  scp -C -i ${ORDERER_HOST_KEY} -r ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}:${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}_${ORDERER_PORT}
  verifyResult $? "Unable to copy ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT} to ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}"


  #create orderer config file
	/bin/cp compose/docker/orderer_template.yaml ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml

  ORDERER_HOME_ESCAPED=`echo ${ORDERER_HOME_REMOTE} | sed -e 's/\//\\\\\//g'`
  sed -i "s/\${ORDERER_PORT}/${ORDERER_PORT}/g"             ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_MSPID}/${ORDERER_MSPID}/g"           ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_HOME}/${ORDERER_HOME_ESCAPED}/g"     ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_ADMIN_PORT}/${ORDERER_ADMIN_PORT}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_OPS_PORT}/${ORDERER_OPS_PORT}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${ORDERER_HOST}/${ORDERER_HOST}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."
  sed -i "s/\${DOMAIN}/${DOMAIN}/g" ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml 
  verifyResult $? "Unable to update orderer.yaml."


  #transfer the config file to orderer machine
  scp -C -i ${ORDERER_HOST_KEY} ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}:${NETWORK_HOME}/${ORDERER_NODE_FOLDER}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml
  verifyResult $? "Unable to copy orderer.yaml to ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN}"

  successln "Successfully created orderer docker compose for port ${ORDERER_PORT} on node $1"
}

function DockerContainerId()
{
  echo $(ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} docker ps -f name=${ORDERER_HOST}.${DOMAIN}_${ORDERER_PORT} -q)
}

function startDockerContainer()
{
  infoln "Starting docker conyainer ..."
  ssh -i ${ORDERER_HOST_KEY} ${ORDERER_HOST_USER}@${ORDERER_HOST}.${DOMAIN} "docker-compose -f ${ORDERER_NODE_HOME}/orderers/${ORDERER_HOST}_${ORDERER_PORT}/orderer.yaml up -d"
}

setOrdererDockerNode $PHYSICAL_SVR_NODE $ORDERER_PORT

CONTAINER_ID=$(DockerContainerId)

if [[ "$CONTAINER_ID" != "" ]];then
  fatalln "There is already a container(id=$CONTAINER_ID) setup, please chose different port or use different orderer node or remove this container."
fi

createDockerCompose 

startDockerContainer

CONTAINER_ID=$(DockerContainerId)
if [[ "$CONTAINER_ID" != "" ]];then
  successln "Container(id=$CONTAINER_ID) is setup and started."
else
  fatalln "Failed to setup and start container"
fi


