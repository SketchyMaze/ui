#!/bin/bash

# Common build script for WASM examples.
# Working directory should be inside a wasm example folder.
# This script creates a "./wasm-tmp" folder, copies the wasm-common files
# into it, copies your .wasm binary in, and runs server.go pointed to that
# directory.
#
# Read: do not run this script yourself. Run it via a `make wasm-serve` command
# in one of the other example directories.
#
# This probably works best on Linux-like systems only.

if [[ ! -f "../wasm-common/serve.sh" ]]; then
	echo Run this script via "make wasm-serve" from a ui/eg example folder.
	exit 1
fi

if [[ -d "./wasm-tmp" ]]; then
	echo Cleaning ./wasm-tmp folder.
	rm -rf ./wasm-tmp
fi

mkdir ./wasm-tmp
cp ../wasm-common/{wasm_exec.js,index.html} ./wasm-tmp/
cp ../DejaVuSans.ttf ./wasm-tmp/
cp *.wasm ./wasm-tmp/app.wasm
cd wasm-tmp/
go run ../../wasm-common/server.go
