#!/bin/bash

set -e

pushd `dirname $0` > /dev/null
scriptdir=`pwd -P`
popd > /dev/null
basedir=`dirname $scriptdir`

echo ""
echo "BUILDING child:"
cd "$basedir/services/child"
sudo docker build -t coopernurse/swarm-demo-child .

echo ""
echo "Pushing docker images to docker hub"
sudo docker push coopernurse/swarm-demo-child

echo ""
echo "Updating docker services"
sudo docker service update --detach=false \
     --env-add '' \
     --update-order=start-first \
     --health-cmd='wget -q -O - http://localhost:9000/env || exit 1' \
     --health-start-period=1s \
     --image docker.io/coopernurse/swarm-demo-child swarm-demo-child

echo ""
echo "DONE"

