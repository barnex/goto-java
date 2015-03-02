#! /bin/bash
go install github.com/barnex/goto-java/cmd/goto-java >> /dev/null 2> /dev/null || (goimports -w *.go &&  go install github.com/barnex/goto-java/cmd/goto-java)
