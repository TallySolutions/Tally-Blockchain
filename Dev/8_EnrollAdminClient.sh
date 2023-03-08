#!/bin/bash

#Usage: 8_EnrollAdminClient.sh 

. ./SetGlobalVariables.sh 1

/bin/rm -rf ${TALLY_HOME}/admin_client

/bin/mkdir -p ${TALLY_HOME}/admin_client

export FABRIC_CA_CLIENT_HOME=${TALLY_HOME}/admin_client
echo "Enrolling Admin Client User"
set -x
fabric-ca-client enroll -u https://${TLS_ADMIN_USER}:${TLS_ADMIN_PASSWORD}@${CA_HOST}.${DOMAIN}:${TLS_CA_PORT} --caname ${TLS_CA_NAME} --enrollment.profile tls --csr.hosts ${CA_HOST}.${DOMAIN} --tls.certfiles "${TLS_CA_HOME}/ca-cert.pem"
{ set +x; } 2>/dev/null

cp ${TALLY_HOME}/admin_client/msp/signcerts/cert.pem ${TALLY_HOME}/admin_client/client-tls-cert.pem
cp ${TALLY_HOME}/admin_client/msp/keystore/* ${TALLY_HOME}/admin_client/client-tls-key.pem

/bin/rm -rf ${TALLY_HOME}/admin_client/msp
