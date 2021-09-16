#!/bin/bash

set -euo pipefail

if ! which CompileDaemon &> /dev/null ; then
	echo "Installing CompileDaemon"
	go get github.com/githubnemo/CompileDaemon
	go install github.com/githubnemo/CompileDaemon
fi

CompileDaemon -build="go build -o mock-server main.go" -command="./mock-server"
