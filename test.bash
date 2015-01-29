#! /bin/bash

set -e

go install
javac Unsigned.java

(cd test && ./test.bash)
