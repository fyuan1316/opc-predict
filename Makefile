BinName ?= opc-predict-server

build-linux: fmt vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/${BinName}-linux main.go
build: fmt vet
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -a -o bin/${BinName} main.go

clean:
	rm -rf bin

vet:
	go vet ./...

fmt:
	go fmt ./...