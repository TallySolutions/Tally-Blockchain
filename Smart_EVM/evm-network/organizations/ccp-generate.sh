#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${ORGMSP}/${1^}MSP/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ${SCRIPTDIR}/organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${ORGMSP}/${1^}MSP/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ${SCRIPTDIR}/organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORG=electorial
P0PORT=7051
CAPORT=7054
PEERPEM=organizations/peerOrganizations/electorial.boatinthesea.com/tlsca/tlsca.electorial.boatinthesea.com-cert.pem
CAPEM=organizations/peerOrganizations/electorial.boatinthesea.com/ca/ca.electorial.boatinthesea.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/electorial.boatinthesea.com/connection-electorial.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/electorial.boatinthesea.com/connection-electorial.yaml

ORG=voter
P0PORT=9051
CAPORT=8054
PEERPEM=organizations/peerOrganizations/voter.boatinthesea.com/tlsca/tlsca.voter.boatinthesea.com-cert.pem
CAPEM=organizations/peerOrganizations/voter.boatinthesea.com/ca/ca.voter.boatinthesea.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/voter.boatinthesea.com/connection-voter.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/voter.boatinthesea.com/connection-voter.yaml
