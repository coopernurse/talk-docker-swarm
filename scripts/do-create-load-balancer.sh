#!/bin/bash

set -e

region="$DO_REGION"

if [[ -z "$region" ]]; then
  region="sfo2"
fi

if [[ -z "$PROJECT" ]]; then
  echo "PROJECT env var must be set"
  exit 1
fi

name="$PROJECT-lb"
tag="$PROJECT-worker"

droplet_ids=$(doctl compute droplet list --tag-name "$tag" | \
  grep -v "Public IPv4" | \
  perl -ne 'print +(split /\s+/)[0], "," ;' | \
  sed 's/.$//')

echo "Creating load balancer: $name in region: $region with droplet ids: $droplet_ids"

doctl compute load-balancer create \
      --name "$name" \
      --droplet-ids "$droplet_ids" \
      --region "$region" \
      --forwarding-rules 'entry_protocol:tcp,entry_port:80,target_protocol:tcp,target_port:80' \
      --health-check 'protocol:http,port:80,path:/,check_interval_seconds:10'
