#!/bin/bash

set -e

pushd `dirname $0` > /dev/null
scriptdir=`pwd -P`
popd > /dev/null
basedir=`dirname $scriptdir`

for name in ui counter clock
do
  image="coopernurse/swarm-demo-$name"
  cd "$basedir/services/$name"
  echo "BUILDING $image"
  sudo docker build -t $image .
  echo "Pushing $image to docker hub"
  sudo docker push $image
done
