all:
	goimports -w *.go
	go install github.com/barnex/goto-java/cmd/goto-java
