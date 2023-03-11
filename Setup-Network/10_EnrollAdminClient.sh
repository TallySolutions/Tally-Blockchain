#!/bin/bash

#Usage: 10_EnrollAdminClient.sh 

. ./SetEnv.sh 
. ./SetCANode.sh 

setCANode 1

/bin/rm -rf ${TALLY_HOME}/admin_client

/bin/mkdir -p ${TALLY_HOME}/admin_client

export FABRIC_CA_CLIENT_HOME=${TALLY_HOME}/admin_client
infoln "Enrolling Admin Client User"
fabric-ca-client enroll -u https://${TLS_ADMIN_USER}:${TLS_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT} --caname ${TLS_CA_NAME} --enrollment.profile tls --csr.hosts ${CA_HOST}.${DOMAIN} --tls.certfiles "${TLS_CA_HOME}/ca-cert.pem"
verifyResult $? "Unable to enroll client admin."


cp ${TALLY_HOME}/admin_client/msp/signcerts/cert.pem ${TALLY_HOME}/admin_client/client-tls-cert.pem
verifyResult $? "Unable to copy tls certificates."
cp ${TALLY_HOME}/admin_client/msp/keystore/* ${TALLY_HOME}/admin_client/client-tls-key.pem
verifyResult $? "Unable to copy tls certificates."

/bin/rm -rf ${TALLY_HOME}/admin_client/msp

successln "Successfully enrolled admin client"