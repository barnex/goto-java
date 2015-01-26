#! /bin/bash

# test script for a single source file files with main()
# runs go and java sources and compares output

set -e

f=$@
basename=$(basename -s ".go" $f)
gofile=$basename.go
jfile=$basename.java
rm -f $gofile.txt $jfile.txt

go run $f >> $gofile.txt 2>> $gofile.txt || (echo "FAIL go run"; exit 1)
goto-java $gofile 2> $gofile.log || (echo "FAIL transpile"; exit 1)
javac $jfile 2> $jfile.log || (echo "FAIL javac"; exit 1)
(cd .. && java main.$basename) >> $basename.java.txt 2> $jfile.txt || (echo "FAIL java"; exit 1)
diff $basename.go.txt $basename.java.txt > $basename.diff || (echo "FAIL diff"; exit 1)
echo "OK"
