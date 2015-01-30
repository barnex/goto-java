#! /bin/bash

# extensive test script

fail=0


printf "%-20s" FILE
for flag in renameall foldconst parens "--"; do
	printf "%-18s" $flag
done
echo

for file in *.go; do
	printf "%-20s" $file 
	for flag in renameall foldconst parens ""; do
		./testone.bash $file --$flag || fail=1
		echo -n -e "\t" 
	done;
	echo 
done;

exit $fail
