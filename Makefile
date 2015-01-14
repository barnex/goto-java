all:
	goimports -w *.go
	go install
