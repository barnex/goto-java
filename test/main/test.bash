#! /bin/bash

# test script for source files with main()
# runs go and java sources and compares output

./clean.bash
javac Unsigned.java || exit 2

failed=0

for f in *.go; do
	echo -e -n $f '\t '
	./testone.bash $f || (( failed++ ))
done;

echo $failed of $(ls *.go | wc -l) failed
exit $failed
