#!/bin/bash

#Usage: 1_SetupCAServers.sh 

#pass default node as 1, node is not used in this script
. ./SetGlobalVariables.sh 1

#TLS

function SetupServer()
{

	CA_HOME=$1
	CA_NAME=$2
	CA_PORT=$3
	CA_USER=$4
	CA_PASSWORD=$5
	CA_OPS_PORT=$6

	echo Creating Fabric CA Server : ${TLS_CA_NAME}

	mkdir -p ${CA_HOME}

	#Create server config file


	/bin/rm -r ${CA_HOME}/*

	/bin/cp fabric-ca-server-config-template.yaml ${CA_HOME}/fabric-ca-server-config.yaml


	sed -e "s/\${CA_NAME}/${CA_NAME}/g" ${CA_HOME}/fabric-ca-server-config.yaml > ${CA_HOME}/fabric-ca-server-config.yaml.1
	sed -e "s/\${CA_PORT}/${CA_PORT}/g" ${CA_HOME}/fabric-ca-server-config.yaml.1 > ${CA_HOME}/fabric-ca-server-config.yaml.2
	sed -e "s/\${CA_USER}/${CA_USER}/g" ${CA_HOME}/fabric-ca-server-config.yaml.2 > ${CA_HOME}/fabric-ca-server-config.yaml.3
	sed -e "s/\${CA_PASSWORD}/${CA_PASSWORD}/g" ${CA_HOME}/fabric-ca-server-config.yaml.3 > ${CA_HOME}/fabric-ca-server-config.yaml.4
	sed -e "s/\${CA_OPS_PORT}/${CA_OPS_PORT}/g" ${CA_HOME}/fabric-ca-server-config.yaml.4 > ${CA_HOME}/fabric-ca-server-config.yaml.5

	/bin/rm ${CA_HOME}/fabric-ca-server-config.yaml ${CA_HOME}/fabric-ca-server-config.yaml.1 ${CA_HOME}/fabric-ca-server-config.yaml.2 ${CA_HOME}/fabric-ca-server-config.yaml.3 ${CA_HOME}/fabric-ca-server-config.yaml.4

	/bin/mv ${CA_HOME}/fabric-ca-server-config.yaml.5 ${CA_HOME}/fabric-ca-server-config.yaml

	fabric-ca-server init --home ${CA_HOME}

	mkdir -p ${CA_HOME}/client

	/bin/cp fabric-ca-client-config-template.yaml ${CA_HOME}/client/fabric-ca-client-config.yaml


	sed -e "s/\${CA_NAME}/${CA_NAME}/g" ${CA_HOME}/client/fabric-ca-client-config.yaml > ${CA_HOME}/client/fabric-ca-client-config.yaml.1
	sed -e "s/\${CA_HOST}/${CA_HOST}/g" ${CA_HOME}/client/fabric-ca-client-config.yaml.1 > ${CA_HOME}/client/fabric-ca-client-config.yaml.2
	sed -e "s/\${DOMAIN}/${DOMAIN}/g" ${CA_HOME}/client/fabric-ca-client-config.yaml.2 > ${CA_HOME}/client/fabric-ca-client-config.yaml.3
	sed -e "s/\${CA_PORT}/${CA_PORT}/g" ${CA_HOME}/client/fabric-ca-client-config.yaml.3 > ${CA_HOME}/client/fabric-ca-client-config.yaml.4

	/bin/rm ${CA_HOME}/client/fabric-ca-client-config.yaml ${CA_HOME}/client/fabric-ca-client-config.yaml.1 ${CA_HOME}/client/fabric-ca-client-config.yaml.2 ${CA_HOME}/client/fabric-ca-client-config.yaml.3 

	/bin/mv ${CA_HOME}/client/fabric-ca-client-config.yaml.4 ${CA_HOME}/client/fabric-ca-client-config.yaml

	echo ========================================

}

#TLS

SetupServer ${TLS_CA_HOME} ${TLS_CA_NAME} ${TLS_CA_PORT} ${TLS_CA_USER} ${TLS_CA_PASSWORD} ${TLS_CA_OPS_PORT}

#Tally

SetupServer ${TALLY_CA_HOME} ${TALLY_CA_NAME} ${TALLY_CA_PORT} ${TALLY_CA_USER} ${TALLY_CA_PASSWORD} ${TALLY_CA_OPS_PORT}

#Orderer

SetupServer ${ORDERER_CA_HOME} ${ORDERER_CA_NAME} ${ORDERER_CA_PORT} ${ORDERER_CA_USER} ${ORDERER_CA_PASSWORD} ${ORDERER_CA_OPS_PORT}


