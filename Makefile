BinName ?= predict-server

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/${BinName}-linux main.go
build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -a -o bin/${BinName} main.go

clean:
	rm -rf bin