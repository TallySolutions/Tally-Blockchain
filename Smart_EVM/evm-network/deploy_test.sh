#!/bin/bash
./network.sh down
./network.sh up
./network.sh createChannel -c election1
./network.sh deployCC -c election1 -ccn evm -ccl go -ccp /home/ubuntu/src/SmartEVM/SmartEVM/chaincode
