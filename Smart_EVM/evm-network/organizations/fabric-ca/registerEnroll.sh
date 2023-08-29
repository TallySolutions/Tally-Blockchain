#!/bin/bash

function createElectorial() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/electorial.boatinthesea.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-electorial --tls.certfiles "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-electorial.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-electorial.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-electorial.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-electorial.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy electorial's CA cert to electorial's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem" "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/msp/tlscacerts/ca.crt"

  # Copy electorial's CA cert to electorial's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem" "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/tlsca/tlsca.electorial.boatinthesea.com-cert.pem"

  # Copy electorial's CA cert to electorial's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/ca"
  cp "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem" "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/ca/ca.electorial.boatinthesea.com-cert.pem"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-electorial --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-electorial --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-electorial --id.name electorialadmin --id.secret electorialadminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-electorial -M "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/msp" --csr.hosts peer0.electorial.boatinthesea.com --tls.certfiles "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-electorial -M "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/tls" --enrollment.profile tls --csr.hosts peer0.electorial.boatinthesea.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/peers/peer0.electorial.boatinthesea.com/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-electorial -M "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/users/User1@electorial.boatinthesea.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/users/User1@electorial.boatinthesea.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://electorialadmin:electorialadminpw@localhost:7054 --caname ca-electorial -M "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/users/Admin@electorial.boatinthesea.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/electorial/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/electorial.boatinthesea.com/users/Admin@electorial.boatinthesea.com/msp/config.yaml"
}

function createVoter() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/voter.boatinthesea.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca-voter --tls.certfiles "${PWD}/organizations/fabric-ca/voter/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-voter.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-voter.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-voter.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-voter.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy voter's CA cert to voter's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/voter/ca-cert.pem" "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/msp/tlscacerts/ca.crt"

  # Copy voter's CA cert to voter's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/voter/ca-cert.pem" "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/tlsca/tlsca.voter.boatinthesea.com-cert.pem"

  # Copy voter's CA cert to voter's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/ca"
  cp "${PWD}/organizations/fabric-ca/voter/ca-cert.pem" "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/ca/ca.voter.boatinthesea.com-cert.pem"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-voter --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/voter/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-voter --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/voter/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-voter --id.name voteradmin --id.secret voteradminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/voter/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-voter -M "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/msp" --csr.hosts peer0.voter.boatinthesea.com --tls.certfiles "${PWD}/organizations/fabric-ca/voter/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-voter -M "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/tls" --enrollment.profile tls --csr.hosts peer0.voter.boatinthesea.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/voter/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/peers/peer0.voter.boatinthesea.com/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca-voter -M "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/users/User1@voter.boatinthesea.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/voter/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/users/User1@voter.boatinthesea.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://voteradmin:voteradminpw@localhost:8054 --caname ca-voter -M "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/users/Admin@voter.boatinthesea.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/voter/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/voter.boatinthesea.com/users/Admin@voter.boatinthesea.com/msp/config.yaml"
}

function createOrderer() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/boatinthesea.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/boatinthesea.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/ordererOrganizations/boatinthesea.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy orderer org's CA cert to orderer org's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/boatinthesea.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/boatinthesea.com/msp/tlscacerts/tlsca.boatinthesea.com-cert.pem"

  # Copy orderer org's CA cert to orderer org's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/ordererOrganizations/boatinthesea.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/boatinthesea.com/tlsca/tlsca.boatinthesea.com-cert.pem"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/msp" --csr.hosts orderer.boatinthesea.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/boatinthesea.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls" --enrollment.profile tls --csr.hosts orderer.boatinthesea.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/server.key"

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/boatinthesea.com/orderers/orderer.boatinthesea.com/msp/tlscacerts/tlsca.boatinthesea.com-cert.pem"

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/boatinthesea.com/users/Admin@boatinthesea.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/boatinthesea.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/boatinthesea.com/users/Admin@boatinthesea.com/msp/config.yaml"
}