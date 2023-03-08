#!/bin/bash

function EnrollUser()
{
        CA_HOME=$1
        URL=$2
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        echo "Enrolling user for ${CA_HOME} [${TYPE}]"
        set -x
        fabric-ca-client enroll -u ${URL} --tls.certfiles "${CA_HOME}/ca-cert.pem"
        { set +x; } 2>/dev/null

}
function RegisterUser()
{
        CA_HOME=$1
        USER=$2
        PASSWORD=$3
        TYPE=$4
        export FABRIC_CA_CLIENT_HOME=${CA_HOME}/client
        echo "Registering user ${USER} for ${CA_HOME} [${TYPE}]"
        set -x
        fabric-ca-client register --id.name ${USER} --id.secret ${PASSWORD} --id.type ${TYPE} --tls.certfiles "${CA_HOME}/ca-cert.pem"
        { set +x; } 2>/dev/null

}


