#!/bin/bash

set -ex
# resolve the directory containing this script
export SCRIPT_DIR=`dirname $0`
pushd $SCRIPT_DIR
export SCRIPT_DIR=$PWD
echo $SCRIPT_DIR
popd

# navigate to the directory containing the schema
pushd $SCRIPT_DIR/../api

# run gqlgen to generate the necessary code based on schema.graphql
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen generate

# navigate back to the root directory
popd
