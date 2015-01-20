#! /bin/bash

set -e

go install

(cd test && ./test.bash)
