#!/bin/ssh
cd golang_ssp

export GOBIN=/usr/local/bin

go install -x

mv /usr/local/bin/golang_ssp /usr/local/bin/ssp