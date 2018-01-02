.PHONY: static
all: static

BINARY=wegigo

dev:
	CGO_ENABLED=0 go build -ldflags="-X deploy.MODE dev" -ldflags -s -o wegigo
static:
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u -v github.com/jteeuwen/go-bindata
	dep ensure
	CGO_ENABLED=0 go build -ldflags -s -o $(BINARY)
	#CGO_ENABLED=0 go build -buildmode=plugin  -o plugins/mod_deploy.so pkg/deploy/module.go
shared:
	CGO_ENABLED=1  go build -ldflags -s -o wegigo
