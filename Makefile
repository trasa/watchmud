all:
	go generate ./...
	go build ./...
	go test ./...
	go fmt ./...
	# build the watchmud exe
	go build .

fmt:
	go fmt ./...

test:
	go generate ./...
	go test ./...
