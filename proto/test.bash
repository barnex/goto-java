#! /bin/bash

go build proto.go
javac Proto.java

rm -f go.out java.out
./proto 2> go.out >> go.out
java Proto 2>> java.out >> java.out
diff go.out java.out
