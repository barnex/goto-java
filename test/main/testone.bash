#! /bin/bash

# test script for a single source file files with main()
# runs go and java sources and compares output

set -e

f=$@
basename=$(basename -s ".go" $f)
gofile=$basename.go
jfile=$basename.java
rm -f $gofile.out $jfile.out $gofile.log $jfile.log $basename.diff $basename.class $jfile

go run $f >> $gofile.out 2>> $gofile.out || (echo "FAIL go run"; exit 1)
goto-java $gofile 2> $gofile.log || (echo "FAIL transpile"; exit 1)
javac $jfile 2> $jfile.log || (echo "FAIL javac"; exit 1)
(cd .. && java main.$basename) >> $jfile.out 2>> $jfile.out || (echo "FAIL java"; exit 1)
diff $gofile.out $jfile.out > $basename.diff || (echo "FAIL diff"; exit 1)
echo "OK"
