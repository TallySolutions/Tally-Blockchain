#!/bin/bash

. ./Utils.sh

#Change this to desired target network setup, Dev: SetEnv-Dev.sh, QA: SetEnv-QA.sh, Prod: SetEnv-Prod.sh and so on...
. ./SetEnv-Dev.sh

. ./SetOrdererNode.sh
. ./SetPeerNode.sh

infoln "Setting up environment ..."
setEnv