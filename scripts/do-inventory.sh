#!/bin/bash

set -e

tag="$1"
invuser="$2"

if [[ -z "$tag" ]]; then
  echo "Usage: do-inventory.sh <tag>"
  exit 1
fi

export invuser

doctl compute droplet list --tag-name "$tag" | \
  grep -v "Public IPv4" | \
  perl -ne '$u=$ENV{"invuser"}; @s = split/\s+/; print "$u\@$s[2]\n"'
