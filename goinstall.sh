#!/bin/ssh
cd golang_ssp

export GOBIN=/usr/local/bin

go install -x

cp /usr/local/bin/golang_ssp /usr/local/bin/ssp
cp /usr/local/bin/golang_ssp /usr/local/bin/ssftp