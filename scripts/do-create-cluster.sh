#!/bin/bash

set -e

project="$1"
sshkey="$2"
region="$DO_REGION"

if [[ -z "$region" ]]; then
  region="sfo2"
fi

if [[ -z "$project" || -z "$sshkey" ]]; then
  echo "Usage: do-create-cluster.sh <project-name> <ssh-key-id>"
  exit 1
fi

echo "Creating droplets for project: $project in region: $region with sshkey: $sshkey"

# first node will be designated the manager
tags="$project-worker,$project-manager"

for (( i=1; i<=3; i++ ))
do
  name="$project-$i"
  echo "creating droplet: $name with tags: $tags"
  doctl compute droplet create "$name" \
        --enable-private-networking \
        --image debian-9-x64 \
        --region "$region" \
        --size 512mb \
        --tag-names "$tags" \
        --ssh-keys "$sshkey"

  # other nodes will be designated workers
  tags="$project-worker"
done
