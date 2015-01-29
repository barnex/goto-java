#! /bin/bash

# test script for source files with main()
# runs go and java sources and compares output

./clean.bash
javac Unsigned.java || exit 2

failed=0

for f in *.go; do
	echo -e -n $f '\t '
	./testone.bash $f $1 || (( failed++ ))
	echo
done;

total=$(ls *.go | wc -l)
passed=$(( $total - $failed ))
echo $passed passed, $failed failed of $total total
exit $failed
