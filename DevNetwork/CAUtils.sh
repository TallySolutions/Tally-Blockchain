#!/bin/bash

function checkIfCAServerRunning()
{
   PID=$(checkIfProcessRunning fabric-ca-server $1)
   if [[ $PID -ne -1 ]]
   then
    warnln "Fabric CA Server $1 is already running on directory '$1' [PID: $PID]"
	  read -p "Want to kill the server to proceed? [y|n]" response
	  if [[ "$response" == "y" ]]
	  then
	     kill -9 $PID
		 verifyResult $? "Could not kill process pid: $PID"
		 return 0
	  else
	     return 1
	  fi
    else
      return 0
   fi
}

export -f checkIfCAServerRunning