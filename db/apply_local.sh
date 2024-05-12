#!/bin/sh
SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

atlas () {
  docker run \
    -v $SCRIPT_DIR/migrations:/migrations \
    -v $SCRIPT_DIR/schema.hcl:/schema.hcl \
    --network=host \
    arigaio/atlas "$@"
}

POSTGRESQL_URL=postgres://postgres:example@localhost:5432/tucows-challenge?sslmode=disable
atlas migrate apply -u ${POSTGRESQL_URL}