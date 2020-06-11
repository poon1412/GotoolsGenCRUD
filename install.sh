#!/bin/sh
export GOPATH=$PWD 
export GOBIN=$PWD/bin 

mkdir $PWD/bin 
export PATH=$PATH:$PWD/bin 

if [ -z "$1" ] || [ "$1" = "glide" ]; then
	curl https://glide.sh/get | sh;

	go get github.com/ngdinhtoan/glide-cleanup;
	cd $GOPATH/src/github.com/ngdinhtoan/glide-cleanup;
	go install;

	go get github.com/multiplay/glide-pin;
	cd $GOPATH/src/github.com/multiplay/glide-pin;
	go install;


	go get github.com/sgotti/glide-vc;
	cd $GOPATH/src/github.com/sgotti/glide-vc;
	go install;

	cd $GOPATH/src/godb2struc;
	glide i -v;
	go install;
fi

if [ -z "$1" ] || [ "$1" = "goapigen" ]; then
	cd $GOPATH/src/goapigen;
	# glide i -v;
	go install;

	# cd $GOPATH/src/gocreate;
	# glide i -v;
	# go install;
fi