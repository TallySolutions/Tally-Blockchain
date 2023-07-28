#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

P0PORT=7051
CAPORT=7054
PEERPEM=${HOME}/fabric/tally-network/organizations/peerOrganizations/tally.tallysolutions.com/tlsca/tlsca.tally.tallysolutions.com-cert.pem
CAPEM=${HOME}/fabric/tally-network/organizations/peerOrganizations/tally.tallysolutions.com/ca/ca.tally.tallysolutions.com-cert.pem

echo "$(json_ccp $P0PORT $CAPORT $PEERPEM $CAPEM)" > ${HOME}/fabric/tally-network/organizations/peerOrganizations/tally.tallysolutions.com/connection-json
echo "$(yaml_ccp $P0PORT $CAPORT $PEERPEM $CAPEM)" > ${HOME}/fabric/tally-network/organizations/peerOrganizations/tally.tallysolutions.com/connection-yaml