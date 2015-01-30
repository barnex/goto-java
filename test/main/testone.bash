#! /bin/bash

# test script for a single source file files with main()
# runs go and java sources and compares output

set -e

w=15
f=$1
basename=$(basename -s ".go" $f)
gofile=$basename.go
jfile=$basename.java
rm -f $gofile.out $jfile.out $gofile.log $jfile.log $basename.diff $basename.class $jfile
go run $f >> $gofile.out 2>> $gofile.out || (printf "%-"$w"s" "FAIL go run"; exit 1)
goto-java $2 $gofile  2> $gofile.log || (printf "%-"$w"s" "FAIL transpile"; exit 1) 
javac $jfile 2> $jfile.log || (printf "%-"$w"s" "FAIL javac"; exit 1) 
(cd .. && java main.$basename) >> $jfile.out 2>> $jfile.out || (printf "%-"$w"s" "FAIL java"; exit 1) 
diff $gofile.out $jfile.out > $basename.diff ||(printf "%-"$w"s" "FAIL diff"; exit 1)  
printf "%-"$w"s" "OK" 
