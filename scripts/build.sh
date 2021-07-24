#!/bin/bash

set -euo pipefail

readonly basedir="$(cd "$(dirname "$0")"/.. && pwd)"
readonly cmd_dir="$basedir/cmd/aviator"

readonly aviator_darwin="aviator-darwin-amd64"
readonly aviator_linux="aviator-linux-amd64"
readonly aviator_win="aviator-win"

main(){
	pushd $cmd_dir
	  go_build "darwin" "$aviator_darwin"
	  go_build "linux" "$aviator_linux"
	  go_build "windows" "$aviator_win"
	popd

	get_shasum "$cmd_dir/$aviator_darwin"
}

go_build(){
  echo "building $2..."
  GOOS=$1 go build -o $2
}

get_shasum(){
	shasum -a 256 $1
}

main
