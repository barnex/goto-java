#! /bin/bash

# test script for source files with main()
# runs go and java sources and compares output


./clean.bash
failed=0

for f in *.go; do
	printf "%-20s" $f;
	./testone.bash $f $1 || (( failed++ ))
	echo
done;
wait

total=$(ls *.go | wc -l)
passed=$(( $total - $failed ))
echo $passed passed, $failed failed of $total total
exit $failed
