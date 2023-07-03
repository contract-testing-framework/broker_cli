BINARY_NAME=broker_cli

build:
	GOARCH=arm64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin-arm64 main.go
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin-amd64 main.go
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux-amd64 main.go
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows-amd64 main.go

test:
	go test ./... -count=1
