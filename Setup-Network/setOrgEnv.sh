#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0




# default to using Tally
ORG=${1:-Tally}

# Exit on first error, print all commands.
set -e
set -o pipefail

# Where am I?
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

ORDERER_CA=${DIR}/tally-network/organizations/ordererOrganizations/tally.tallysolutions.com/tlsca/tlsca.tally.tallysolutions.com-cert.pem
PEER0_ORG1_CA=${DIR}/tally-network/organizations/peerOrganizations/tally.tallysolutions.com/tlsca/tlsca.tally.tallysolutions.com-cert.pem
PEER0_ORG2_CA=${DIR}/tally-network/organizations/peerOrganizations/org2.tallysolutions.com/tlsca/tlsca.org2.tallysolutions.com-cert.pem
PEER0_ORG3_CA=${DIR}/tally-network/organizations/peerOrganizations/org3.tallysolutions.com/tlsca/tlsca.org3.tallysolutions.com-cert.pem


if [[ ${ORG,,} == "tally" || ${ORG,,} == "digibank" ]]; then

   CORE_PEER_LOCALMSPID=Tally
   CORE_PEER_MSPCONFIGPATH=${DIR}/tally-network/organizations/peerOrganizations/tally.tallysolutions.com/users/Admin@tally.tallysolutions.com/msp
   CORE_PEER_ADDRESS=localhost:7051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/tally-network/organizations/peerOrganizations/tally.tallysolutions.com/tlsca/tlsca.tally.tallysolutions.com-cert.pem

elif [[ ${ORG,,} == "org2" || ${ORG,,} == "magnetocorp" ]]; then

   CORE_PEER_LOCALMSPID=Org2MSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/tally-network/organizations/peerOrganizations/org2.tallysolutions.com/users/Admin@org2.tallysolutions.com/msp
   CORE_PEER_ADDRESS=localhost:9051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/tally-network/organizations/peerOrganizations/org2.tallysolutions.com/tlsca/tlsca.org2.tallysolutions.com-cert.pem

else
   echo "Unknown \"$ORG\", please choose Tally/Digibank or Org2/Magnetocorp"
   echo "For example to get the environment variables to set upa Org2 shell environment run:  ./setOrgEnv.sh Org2"
   echo
   echo "This can be automated to set them as well with:"
   echo
   echo 'export $(./setOrgEnv.sh Org2 | xargs)'
   exit 1
fi

# output the variables that need to be set
echo "CORE_PEER_TLS_ENABLED=true"
echo "ORDERER_CA=${ORDERER_CA}"
echo "PEER0_ORG1_CA=${PEER0_ORG1_CA}"
echo "PEER0_ORG2_CA=${PEER0_ORG2_CA}"
echo "PEER0_ORG3_CA=${PEER0_ORG3_CA}"

echo "CORE_PEER_MSPCONFIGPATH=${CORE_PEER_MSPCONFIGPATH}"
echo "CORE_PEER_ADDRESS=${CORE_PEER_ADDRESS}"
echo "CORE_PEER_TLS_ROOTCERT_FILE=${CORE_PEER_TLS_ROOTCERT_FILE}"

echo "CORE_PEER_LOCALMSPID=${CORE_PEER_LOCALMSPID}"
