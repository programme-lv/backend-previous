#!/bin/bash

set -ex # fail fast

# resolve the directory containing this script
export SCRIPT_DIR=`dirname $0`
pushd $SCRIPT_DIR
export SCRIPT_DIR=$PWD
echo $SCRIPT_DIR
popd

pushd $SCRIPT_DIR/../internal/database

jet -dsn="postgresql://proglv:proglv@localhost:5432/proglv?sslmode=disable" -path=.

popd