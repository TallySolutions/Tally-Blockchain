#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

# imports
. ${SCRIPTDIR}/scripts/utils.sh

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/boatinthesea.com/tlsca/tlsca.boatinthesea.com-cert.pem
export PEER0_ELECTORIAL_CA=${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/tlsca/tlsca.electorial.boatinthesea.com-cert.pem
export PEER0_VOTER_CA=${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/tlsca/tlsca.voter.boatinthesea.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/server.key

# Set environment variables for the peer org
setGlobals() {
  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  infoln "Using organization ${USING_ORG}"
  if [ $USING_ORG == "electorial" ]; then
    export CORE_PEER_LOCALMSPID="ElectorialMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ELECTORIAL_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/users/Admin@electorial.boatinthesea.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
  elif [ $USING_ORG == "voter" ]; then
    export CORE_PEER_LOCALMSPID="VoterMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_VOTER_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/users/Admin@voter.boatinthesea.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
  else
    errorln "ORG Unknown"
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

# Set environment variables for use in the CLI container
setGlobalsCLI() {
  setGlobals $1

  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  if [ $USING_ORG == "electorial" ]; then
    export CORE_PEER_ADDRESS=peer0.electorial.boatinthesea.com:7051
  elif [ $USING_ORG == "voter" ]; then
    export CORE_PEER_ADDRESS=peer0.voter.boatinthesea.com:9051
  else
    errorln "ORG Unknown"
  fi
}

# parsePeerConnectionParameters $@
# Helper function that sets the peer connection parameters for a chaincode
# operation
parsePeerConnectionParameters() {
  PEER_CONN_PARMS=()
  PEERS=""
  while [ "$#" -gt 0 ]; do
    setGlobals $1
    PEER="peer0.$1"
    ## Set peer addresses
    if [ -z "$PEERS" ]
    then
	PEERS="$PEER"
    else
	PEERS="$PEERS $PEER"
    fi
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" --peerAddresses $CORE_PEER_ADDRESS)
    ## Set path to TLS certificate
    CA=PEER0_${1^^}_CA
    TLSINFO=(--tlsRootCertFiles "${!CA}")
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" "${TLSINFO[@]}")
    # shift by one to get to the next organization
    shift
  done
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}
