
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./_build/btcd-linux

.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ./_build/btcd-darwin

