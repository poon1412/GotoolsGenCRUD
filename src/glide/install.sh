#!/bin/sh
export GOPATH=$PWD/../../.  
export GOBIN=$PWD/../../bin 
curl https://glide.sh/get | sh


go get github.com/ngdinhtoan/glide-cleanup
cd $GOPATH/src/github.com/ngdinhtoan/glide-cleanup
go build

go get github.com/multiplay/glide-pin
cd $GOPATH/src/github.com/multiplay/glide-pin
go build


go get github.com/sgotti/glide-vc
cd $GOPATH/src/github.com/sgotti/glide-vc
go build

