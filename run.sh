#!/bin/sh
export GOPATH=$PWD
export GOBIN=$PWD/bin
 
go fmt $1/... 


if [[ $RES = *"FAIL"* || $RES = *"failed"* ]]; then 
	echo "Stop run project, Because Test Fail." ; 
	exit;
fi

go run src/$1/main.go