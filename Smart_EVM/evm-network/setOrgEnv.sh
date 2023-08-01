#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0




# default to using Electorial
ORG=${1:-Electorial}

# Exit on first error, print all commands.
set -e
set -o pipefail

# Where am I?
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

ORDERER_CA=${DIR}/evm-network/organizations/ordererOrganizations/boatinthesea.com/tlsca/tlsca.boatinthesea.com-cert.pem
PEER0_ELECTORIAL_CA=${DIR}/evm-network/organizations/peerOrganizations/electorial.boatinthesea.com/tlsca/tlsca.electorial.boatinthesea.com-cert.pem
PEER0_VOTER_CA=${DIR}/evm-network/organizations/peerOrganizations/voter.boatinthesea.com/tlsca/tlsca.voter.boatinthesea.com-cert.pem


if [[ ${ORG,,} == "electorial" ]]; then

   CORE_PEER_LOCALMSPID=ElectorialMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/evm-network/organizations/peerOrganizations/electorial.boatinthesea.com/users/Admin@electorial.boatinthesea.com/msp
   CORE_PEER_ADDRESS=localhost:7051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/evm-network/organizations/peerOrganizations/electorial.boatinthesea.com/tlsca/tlsca.electorial.boatinthesea.com-cert.pem

elif [[ ${ORG,,} == "voter" ]]; then

   CORE_PEER_LOCALMSPID=VoterMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/evm-network/organizations/peerOrganizations/voter.boatinthesea.com/users/Admin@voter.boatinthesea.com/msp
   CORE_PEER_ADDRESS=localhost:9051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/evm-network/organizations/peerOrganizations/voter.boatinthesea.com/tlsca/tlsca.voter.boatinthesea.com-cert.pem

else
   echo "Unknown \"$ORG\", please choose Electorial or Voter"
   echo "For example to get the environment variables to set upa Voter shell environment run:  ./setOrgEnv.sh Voter"
   echo
   echo "This can be automated to set them as well with:"
   echo
   echo 'export $(./setOrgEnv.sh Voter | xargs)'
   exit 1
fi

# output the variables that need to be set
echo "CORE_PEER_TLS_ENABLED=true"
echo "ORDERER_CA=${ORDERER_CA}"
echo "PEER0_ELECTORIAL_CA=${PEER0_ELECTORIAL_CA}"
echo "PEER0_VOTER_CA=${PEER0_VOTER_CA}"

echo "CORE_PEER_MSPCONFIGPATH=${CORE_PEER_MSPCONFIGPATH}"
echo "CORE_PEER_ADDRESS=${CORE_PEER_ADDRESS}"
echo "CORE_PEER_TLS_ROOTCERT_FILE=${CORE_PEER_TLS_ROOTCERT_FILE}"

echo "CORE_PEER_LOCALMSPID=${CORE_PEER_LOCALMSPID}"
