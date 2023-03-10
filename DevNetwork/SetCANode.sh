#!/bin/bash

function setCANode()
{

  NODE=`formatNodeNo $1`

  CA_HOST=${CA_HOST_PREFIX}${NODE}

}

export -f setCANode