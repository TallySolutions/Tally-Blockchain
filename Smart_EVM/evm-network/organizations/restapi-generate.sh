#!/bin/bash

#imports
. ${SCRIPTDIR}/scripts/utils.sh

function verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}

function gen_script(){ 

    ORG=$1
    PEER=$2
    PORT=$3
    MSP=$4

    SERVICE_STARUP_DIR=organizations/peerOrganizations/${ORG}.boatinthesea.com/peers/${PEER}.${ORG}.boatinthesea.com/bin

    mkdir -p ${SERVICE_STARUP_DIR}

    SERVICE_STARUP_SCRIPT=${SERVICE_STARUP_DIR}/start.sh

    /bin/cp ${WEB_SERVER_SOURCE}/web ${SERVICE_STARUP_DIR}/web
    verifyResult $? "Web server copy failed."

    echo "#!/bin/bash" > ${SERVICE_STARUP_SCRIPT}
    echo >> ${SERVICE_STARUP_SCRIPT}
    echo "#Start the peer process" >> ${SERVICE_STARUP_SCRIPT}
    echo "peer node start &" >> ${SERVICE_STARUP_SCRIPT}
    echo >> ${SERVICE_STARUP_SCRIPT}
    echo "#Start the web service" >> ${SERVICE_STARUP_SCRIPT}
    echo "/etc/hyperledger/fabric/bin/web ${PEER} ${ORG} ${PORT} boatinthesea.com ${MSP} &" >> ${SERVICE_STARUP_SCRIPT}
    echo >> ${SERVICE_STARUP_SCRIPT}
    echo "# Wait for any process to exit" >> ${SERVICE_STARUP_SCRIPT}
    echo "wait -n" >> ${SERVICE_STARUP_SCRIPT}
    echo >> ${SERVICE_STARUP_SCRIPT}
    echo "# Exit with status of process that exited first" >> ${SERVICE_STARUP_SCRIPT}
    echo "exit $?" >> ${SERVICE_STARUP_SCRIPT}
    verifyResult $? "Start script preparation failed."

    chmod +x ${SERVICE_STARUP_SCRIPT}
    verifyResult $? "Start script preparation failed (error setting execute permission)."

}

WEB_SERVER_SOURCE=${SCRIPTDIR}/../SmartEVM/rest-api
#Build web server
#CDIR=${PWD}
#cd ${WEB_SERVER_SOURCE}
docker run -v ${SCRIPTDIR}/..:/SmartEVM -e GOPATH=/opt/gopath -w /SmartEVM/SmartEVM/rest-api --rm hyperledger/fabric-tools:latest go build  -o web
#CGO_ENABLED=0 go build -o /tmp/web 
res=$?
#cd ${CDIR}
{ set +x; } 2>/dev/null
verifyResult $res "Web server build failed."
successln "Successfully built the web server."


#Generate script for the electorial/peer0
gen_script electorial peer0 7051 ElectorialMSP

#Generate script for the voter/peer0
gen_script voter peer0 9051 VoterMSP


/bin/rm -f ${WEB_SERVER_SOURCE}/web