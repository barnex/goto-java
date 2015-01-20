#! /bin/bash

# test script for source files with main()
# runs go and java sources and compares output

set -e

a3 *.go

javac *.java
rm -f *.txt

for f in *.java; do
	basename=$(basename -s ".java" $f)
	(cd .. && java main.$basename) > $basename.java.txt
	go run $basename.go 2> $basename.go.txt
	diff $basename.go.txt $basename.java.txt
done;

