#!/bin/bash

export GOPATH=`pwd`
#export GOPATH=$(pwd):$(pwd)/../common
export GOBIN=`pwd`/bin
export GOOS="linux"


echo "GOOS:$GOOS"
echo "GOBIN:$GOBIN"
echo "GOPATH:$GOPATH"
