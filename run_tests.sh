#!/bin/sh

cd ~/go/src/github.com/zbo14/envoke/crypto
go test -v

cd ~/go/src/github.com/zbo14/envoke/bigchain 
go test -v

cd ~/go/src/github.com/zbo14/envoke/api
go test -v