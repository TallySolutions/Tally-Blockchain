#!/bin/bash

#Usage: 1_SetupCAServers.sh 

. ./SetEnv.sh 
. ./CAUtils.sh
. ./SetCANode.sh 

setCANode 1

function SetupServer()
{

	CA_HOME=$1
	CA_NAME=$2
	CA_PORT=$3
	CA_USER=$4
	CA_PASSWORD=$5
	CA_OPS_PORT=$6

	infoln Creating Fabric CA Server : ${CA_NAME}

	mkdir -p ${CA_HOME}

    checkIfCAServerRunning $CA_HOME
	verifyResult $? "Fabric CA Server $CA_NAME already setup at $CA_HOME is running, please stop them first."

	#Create server config file

    infoln "Cleaning up $CA_HOME ..."
	/bin/rm -r ${CA_HOME}/*
	verifyResult $? "Could not cleanup $CA_HOME."

    infoln "Creating ${CA_HOME}/fabric-ca-server-config.yaml ..."
	/bin/cp fabric-ca-server-config-template.yaml ${CA_HOME}/fabric-ca-server-config.yaml
	verifyResult $? "Could not create $CA_HOME/fabric-ca-server-config.yaml."

	sed -i "s/\${CA_HOST}/${CA_HOST}/g" ${CA_HOME}/fabric-ca-server-config.yaml 
	verifyResult $? "Error updating $CA_HOME/fabric-ca-server-config.yaml."
	sed -i "s/\${DOMAIN}/${DOMAIN}/g" ${CA_HOME}/fabric-ca-server-config.yaml 
	verifyResult $? "Error updating $CA_HOME/fabric-ca-server-config.yaml."
	sed -i "s/\${CA_NAME}/${CA_NAME}/g" ${CA_HOME}/fabric-ca-server-config.yaml 
	verifyResult $? "Error updating $CA_HOME/fabric-ca-server-config.yaml."
	sed -i "s/\${CA_PORT}/${CA_PORT}/g" ${CA_HOME}/fabric-ca-server-config.yaml
	verifyResult $? "Error updating $CA_HOME/fabric-ca-server-config.yaml."
	sed -i "s/\${CA_USER}/${CA_USER}/g" ${CA_HOME}/fabric-ca-server-config.yaml
	verifyResult $? "Error updating $CA_HOME/fabric-ca-server-config.yaml."
	sed -i "s/\${CA_PASSWORD}/${CA_PASSWORD}/g" ${CA_HOME}/fabric-ca-server-config.yaml
	verifyResult $? "Error updating $CA_HOME/fabric-ca-server-config.yaml."
	sed -i "s/\${CA_OPS_PORT}/${CA_OPS_PORT}/g" ${CA_HOME}/fabric-ca-server-config.yaml
	verifyResult $? "Error updating $CA_HOME/fabric-ca-server-config.yaml."


    infoln "Initializing $CA_NAME CA Server ..."
	fabric-ca-server init --home ${CA_HOME}
	verifyResult $? "Error initializing $CA_NAME CA Server."

    infoln "Creating client config for $CA_NAME ..."
	mkdir -p ${CA_HOME}/client

	/bin/cp fabric-ca-client-config-template.yaml ${CA_HOME}/client/fabric-ca-client-config.yaml
	verifyResult $? "Could not create $CA_HOME/client/fabric-ca-client-config.yaml."


	sed -i "s/\${CA_NAME}/${CA_NAME}/g" ${CA_HOME}/client/fabric-ca-client-config.yaml 
	verifyResult $? "Error updating $CA_HOME/client/fabric-ca-client-config.yaml."
	sed -i "s/\${CA_HOST}/${CA_HOST}/g" ${CA_HOME}/client/fabric-ca-client-config.yaml
	verifyResult $? "Error updating $CA_HOME/client/fabric-ca-client-config.yaml."
	sed -i "s/\${DOMAIN}/${DOMAIN}/g" ${CA_HOME}/client/fabric-ca-client-config.yaml
	verifyResult $? "Error updating $CA_HOME/client/fabric-ca-client-config.yaml."
	sed -i "s/\${CA_PORT}/${CA_PORT}/g" ${CA_HOME}/client/fabric-ca-client-config.yaml
	verifyResult $? "Error updating $CA_HOME/client/fabric-ca-client-config.yaml."

	successln "$CA_NAME CA Server created successfully."

}

#TLS

SetupServer ${TLS_CA_HOME} ${TLS_CA_NAME} ${TLS_CA_PORT} ${TLS_CA_USER} ${TLS_CA_PASSWORD} ${TLS_CA_OPS_PORT}

#Tally

SetupServer ${TALLY_CA_HOME} ${TALLY_CA_NAME} ${TALLY_CA_PORT} ${TALLY_CA_USER} ${TALLY_CA_PASSWORD} ${TALLY_CA_OPS_PORT}

#Orderer

SetupServer ${ORDERER_CA_HOME} ${ORDERER_CA_NAME} ${ORDERER_CA_PORT} ${ORDERER_CA_USER} ${ORDERER_CA_PASSWORD} ${ORDERER_CA_OPS_PORT}


