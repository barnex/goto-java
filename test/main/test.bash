#! /bin/bash

# test script for source files with main()
# runs go and java sources and compares output

./clean.bash

rm -f *.java

exitstat=0

for f in *.go; do
	echo -e -n $f '\t '
	./testone.bash $f || exitstat=1
done;

exit $exitstat
