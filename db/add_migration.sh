#!/bin/sh
set -x
SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

atlas () {
  docker run \
    -v $SCRIPT_DIR/migrations:/migrations \
    -v $SCRIPT_DIR/schema.hcl:/schema.hcl \
    --network=host \
    arigaio/atlas "$@"
}

atlas migrate diff $1 --dir="file://migrations" --to="file://schema.hcl" --dev-url="postgres://user1:pw1@localhost:5433/atlas-migration-compare?sslmode=disable"