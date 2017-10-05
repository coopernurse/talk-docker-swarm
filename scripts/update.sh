#!/bin/bash

set -e

echo ""
echo "Updating docker services"
sudo docker service update --detach=false \
     --image docker.io/coopernurse/swarm-demo-parent swarm-demo-parent
sudo docker service update --detach=false \
     --image docker.io/coopernurse/swarm-demo-child swarm-demo-child

echo ""
echo "DONE"

