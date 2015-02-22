#! /bin/bash

set -e

mkdir -p build

for f in *.go; do
	echo go build -o build/$(basename -s ".go" build/$f) $f
	go build -o build/$(basename -s ".go" build/$f) $f
done
