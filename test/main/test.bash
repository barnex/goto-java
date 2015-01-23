#! /bin/bash

# test script for source files with main()
# runs go and java sources and compares output

./clean.bash

for par in 1 0; do

echo goto-java -print=false -parens=$par *.go
goto-java -print=false -parens=$par *.go
javac *.java

for f in *.java; do
	echo $f
	basename=$(basename -s ".java" $f)
	(cd .. && java main.$basename) > $basename.java.txt
	go run $basename.go 2> $basename.go.txt
	diff $basename.go.txt $basename.java.txt || exit 1;
done;

done;
