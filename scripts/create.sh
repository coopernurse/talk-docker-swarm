#!/bin/bash

set -e

sudo docker network create -d overlay net1
sudo docker service create --network net1 --replicas 1 \
     --publish mode=host,target=9000,published=9000 \
     --name swarm-demo-parent docker.io/coopernurse/swarm-demo-parent
sudo docker service create --network net1 --replicas 1 \
     --publish 9000 \
     --name swarm-demo-child docker.io/coopernurse/swarm-demo-child
