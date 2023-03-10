#!/bin/bash

C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

# println echos string
function println() {
  echo -e "$*"
}

# errorln echos i red color
function errorln() {
  println "${C_RED}${*}${C_RESET}"
}

# successln echos in green color
function successln() {
  println "${C_GREEN}${*}${C_RESET}"
}

# infoln echos in blue color
function infoln() {
  println "${C_BLUE}${*}${C_RESET}"
}

# warnln echos in yellow color
function warnln() {
  println "${C_YELLOW}${*}${C_RESET}"
}

# fatalln echos in red color and exits with fail status
function fatalln() {
  errorln "$*"
  exit 1
}

function formatNodeNo()
{

   NODE=${1:-1}

   re='^[0-9]+$'
   if ! [[ $NODE =~ $re ]] ; then
      fatalln "Usage: Node value should be number!"
   fi

   if [[ ${NODE} -lt 10 ]] ; then
      NODE="0${NODE}"
    fi
   
   echo $NODE

}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}

warnResult() {
  if [ $1 -ne 0 ]; then
    warnln "$2"
  fi
}

function checkIfProcessRunning()
{
	PNAME=$1
	PDIR=$2
	PIDS=$(pgrep -f $PNAME)

	while read PROCESS 
	do 
	    PID=$(echo $PROCESS | cut -f1 -d:)
        WD=$(echo $PROCESS | cut -f2 -d: | sed 's/^[ \t]*//')
        if [[ "$WD" == "$PDIR" ]]
		then 
		  echo $PID
		  return 
		fi        
    done << EOF
     	$(pwdx $PIDS 2> /dev/null)
EOF
	echo -1
}

function checkIfRemoteProcessRunning()
{
	PNAME=$1
	RUSER=$2
  RSVR=$3
  RKEY=$4

	PIDS=$(ssh -i $RKEY $RUSER@$RSVR pgrep -f $PNAME)

	for PID in $PIDS 
	do 
		  echo $PID
		  return 
  done 
	echo -1
}


export -f errorln
export -f successln
export -f infoln
export -f warnln
export -f formatNodeNo
export -f verifyResult
export -f checkIfProcessRunning