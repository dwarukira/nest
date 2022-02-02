#!/bin/bash
set -e

if [ -f .env.local ]; then
  export $(cat .env.local | grep -v '#' | awk '/=/ {print $1}')
fi

make migrate-create -command=$1
