#!/bin/sh

source init.sh

read -p "Enter path to audio file: " path

export PATH_TO_AUDIO_FILE=$path

cd ~/go/src/github.com/zbo14/envoke/crypto
go test -v

cd ~/go/src/github.com/zbo14/envoke/bigchain 
go test -v

cd ~/go/src/github.com/zbo14/envoke/api
go test -v