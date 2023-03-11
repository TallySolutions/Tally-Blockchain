#!/bin/bash

function EnrollUser()
{
        CA_HOME=$1
        URL=$2
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        infoln "Enrolling user for ${CA_HOME} [${TYPE}]"
        set -x
        fabric-ca-client enroll -u ${URL} --tls.certfiles "${CA_HOME}/ca-cert.pem"
        res=$?
        { set +x; } 2>/dev/null
        verifyResult $res "Unable to enroll user for ${CA_HOME} [${TYPE}]. (Is relevant CA server running?)"

}
function RegisterUser()
{
        CA_HOME=$1
        USER=$2
        PASSWORD=$3
        TYPE=$4
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        infoln "Registering user ${USER} for ${CA_HOME} [${TYPE}]"
        set -x
        fabric-ca-client register --id.name ${USER} --id.secret ${PASSWORD} --id.type ${TYPE} --tls.certfiles "${CA_HOME}/ca-cert.pem"
        res=$?
        { set +x; } 2>/dev/null
        verifyResult $res "Unable to register user for ${CA_HOME} [${TYPE}]. (Are you trying to re-register a user or CA server is stopped?)"

}


