all:
	go build ./...
	go test ./...

fmt:
	go fmt ./...