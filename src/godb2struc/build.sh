#!/bin/sh
export GOPATH=$PWD/../../  
export GOBIN=$PWD/../../bin 
# go install
go build
# go build -o ~/Documents/workspace/go-project/tools/godb2struc
# godb2struc -user=devorisma -pass=0rism@** -db=SCB18B-HOME-TEST -host=10.91.2.45 -o=$GOPATH/src/godb2struc/datamodels