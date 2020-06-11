#!/bin/sh
export GOPATH=$PWD
export GOBIN=$PWD/bin
cd $GOPATH/src/$1
if [ ! -f "$GOPATH/src/$1/glide.yaml" ]; then 
	glide create --non-interactive && glide i -v; 
else 
	glide up -v; 
fi