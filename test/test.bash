set -e

(cd .. && go install)
rm -f *.java *.class

(cd main && ./test.bash)

