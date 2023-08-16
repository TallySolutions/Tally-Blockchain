#!/bin/bash
./network.sh down
./network.sh up
./network.sh createChannel -c tally
./network.sh deployCC -c tally -ccn evm -ccl go -ccp /home/ubuntu/Tally-Blockchain/Smart_EVM
 