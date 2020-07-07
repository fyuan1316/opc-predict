build-predict-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/opc-predict cmd/predict/main.go
build-predict:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -a -o bin/opc-predict-darwin cmd/predict/main.go

build-predict-server-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/opc-predict-server cmd/server/main.go

build-predict-server:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -a -o bin/opc-predict-server-darwin cmd/server/main.go

clean:
	rm -rf bin