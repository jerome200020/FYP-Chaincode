#!/bin/bash
function _exit(){
    printf "Exiting:%s\n" "$1"
    exit -1
}

# Exit on first error, print all commands.
set -ev
set -o pipefail

# Where am I?
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export FABRIC_CFG_PATH="${DIR}/../config"

cd "../test-network/"

./network.sh down
./network.sh up createChannel -c test -s couchdb -ca

./network.sh deployCC -c test -ccn cert -ccp ../fyp-chaincode/chaincode -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"

echo "Set up completed"
