#!/bin/bash

set -e pipefail

readonly scriptname=$(basename -- $0)
readonly basedir=$(dirname $scriptname)

readonly aviator_darwin="aviator-darwin-amd64"
readonly aviator_linux="aviator-linux-amd64"
readonly aviator_win="aviator-win"

main(){
	pushd $basedir/cmd/aviator
	  go_build "darwin" "$aviator_darwin"
	  go_build "linux" "$aviator_linux"
	  go_build "windows" "$aviator_win"
	popd

	get_shasum "$basedir/cmd/aviator/$aviator_darwin"
}

go_build(){
  echo "building $2..."
  GOOS=$1 go build -o $2
}

get_shasum(){
	shasum -a 256 $1
}

main
